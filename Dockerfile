FROM gcr.io/distroless/static:nonroot

COPY kube-token-refresher /usr/local/bin/ktr

ENTRYPOINT [ "/usr/local/bin/ktr" ]
