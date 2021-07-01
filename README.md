# kube-token-refresher

![](https://img.shields.io/github/go-mod/go-version/vshn/kube-token-refresher)
[![](https://img.shields.io/github/license/vshn/kube-token-refresher)](https://github.com/vshn/kube-token-refresher/blob/master/LICENSE)

The `kube-token-refresher` is a tool to periodically fetch access token and write them to Kubernetes secrets.
It fetches the access token using the OpenId Connect Client Credentials grant.

This enables systems that expect long lived access token to work with short lived tokens that expire frequently.

## Configuration

The `kube-token-refresher` can be configured using a YAML file, environment variables, or command line flags.

``` yaml
---
## kube-token-refresher configuration

secret:
  # Name of the secret to write to
  name: ''

  # Namespace of the secret to write to
  namespace: ''

  # The key in the specified secret to write the token to
  key: 'token'

# In what interval (in seconds) to fetch a new token and update the secret
# You should count in possible timeouts upon refreshing and also mount update, see
# https://kubernetes.io/docs/concepts/configuration/secret/#mounted-secrets-are-updated-automatically
interval: 500

log:
  # How verbose the logging should be. One of:
  # * debug
  # * info
  # * warn
  level: 'info'

  # What format to log in. One of
  # * text
  # * json
  format: 'text'

# Configures how to connect to the OIDC provider
oidc:
  # The toke endpoint of the OpenId Connect provider
  #
  # Usually in the form of: `https://<domain>/token`
  tokenurl: ''

  # The Client ID
  clientid: ''

  # The Client Secret
  clientsecret: ''
```

The configuration file has to be specified directly with the `--config` flag.

### Environment Variables

All configuration values can be set through environment variables.
The configuration key is translated to a environment variable with the prefix `KTR_` and the key name in all caps.
For nested configuration keys the levels are separated with `_`.


```
# Will set the logLevel to `debug`
export KTR_LOG_LEVEL="debug"

# Will set the OIDC tokenUrl to `https://auth.vshn.net/token`
export KTR_OIDC_TOKENURL="https://auth.vshn.net/token"

```

Environment variables will take precedence over the configuration file.


### Command Line Flags

All configuration values can also be set through command line flags.
The configuration key is directly translated to the flag.
Nested configuration keys are separated with `.`.

```
# Will set the logLevel to `warn`
./kube-token-refresher --log.level warn

# Will set the OIDC tokenUrl to `https://auth.vshn.net/token`
./kube-token-refresher --oidc.tokenurl="https://auth.vshn.net/token"
```


Command line flags will take precedence over both the configuration file and environment variables.


## Deploy

To deploy the `kube-token-refresher` you need OIDC credentials capable of requesting an access token, and Kubernetes credentials to `get` and `update` the specified secret.
If the `kube-token-refresher` is expected to create the specified secret, it will also need permission to `create` secrets.

You can find an example deployment in `deploy/`.
