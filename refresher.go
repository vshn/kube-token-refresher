package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type tokenProvider interface {
	GetToken(context.Context) ([]byte, error)
}

type dummyProvider struct {
	i int
}

func (d *dummyProvider) GetToken(context.Context) ([]byte, error) {
	d.i = d.i + 1
	return []byte(fmt.Sprintf("token %d", d.i)), nil
}

type refresher struct {
	name      string
	namespace string
	key       string

	kclient.Client
	provider tokenProvider
}

func (r refresher) refresh(ctx context.Context) error {
	t, err := r.provider.GetToken(ctx)
	if err != nil {
		return fmt.Errorf("error getting new token: %w", err)
	}

	secret := &corev1.Secret{}
	err = r.Get(ctx, kclient.ObjectKey{
		Name:      r.name,
		Namespace: r.namespace,
	}, secret)
	if err != nil && !kerrors.IsNotFound(err) {
		return fmt.Errorf("error getting secret: %w", err)
	}

	if kerrors.IsNotFound(err) {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      r.name,
				Namespace: r.namespace,
			},
			Data: map[string][]byte{r.key: t},
			Type: corev1.SecretTypeOpaque,
		}
		return r.Create(ctx, secret)
	}
	patch := &corev1.Secret{
		ObjectMeta: secret.ObjectMeta,
		Data:       map[string][]byte{r.key: t},
	}

	return r.Patch(ctx, patch, kclient.StrategicMergeFrom(secret))
}
