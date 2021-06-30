# kube-token-refresher

![](https://img.shields.io/github/go-mod/go-version/vshn/kube-token-refresher)
[![](https://img.shields.io/github/license/vshn/kube-token-refresher)](https://github.com/vshn/kube-token-refresher/blob/master/LICENSE)

The `kube-token-refresher` is a tool to periodically fetch access token and 
write them to Kubernetes secrets. It fetches the access token using the 
OpenId Connect Client Credentials grant.

This enables systems that expect long lived access token to work with short
lived tokens that expire frequently.

## Configuration

The `kube-token-refresher` can be configured using a YAML file, environment 
variables, or command line flags.

``` yaml
---
## kube-token-refresher configuration

# Name of the secret to write to
secretName: ''

# Namespace of the secret to write to
secretNamespace: ''

# The key in the specified secret to write the token to
secretKey: 'token'

# In what interval (in seconds) to fetch a new token and update the secret
refreshInterval: 595

# How verbose the logging should be. One of:
# * debug
# * info
# * warn
logLevel: 'info'

# What format to log in. One of
# * text
# * json
logFormat: 'text'

# Configures how to connect to the OIDC provider
oidc:
  # The toke endpoint of the OpenId Connect provider
  # 
  # Usually in the form of: `https://<domain>/token`
  tokenUrl: ''

  # The Client ID 
  clientID: ''

  # The Client Secret
  clientSecret: ''
```

The configuration file can be located in one of the following paths or be 
specified directly with the `--config` flag.

* `/etc/kube-token-refresher/config.yml`
* `$HOME/.config/kube-token-refresher/config.yml`
* `$HOME/.kube-token-refresher/config.yml`

### Environment Variables

All configuration values can be set through environment variables. The 
configuration key is translated to a environment variable with the prefix
`KTR_` and the key name in all caps. For nested configuration keys the 
levels are separated with `_`.


```
# Will set the logLevel to `debug`
export KTR_LOGLEVEL="debug"

# Will set the OIDC tokenUrl to `https://auth.vshn.net/token`
export KTR_OIDC_TOKENURL="https://auth.vshn.net/token"

```

Environment variables will take precedence over the configuration file.


### Command Line Flags

All configuration values can also be set through command line flags. The 
configuration key is directly translated to the flag. Nested configuration
keys are separated with `.`.

```
# Will set the logLevel to `warn`
./kube-token-refresher --logLevel warn

# Will set the OIDC tokenUrl to `https://auth.vshn.net/token`
./kube-token-refresher --oidc.tokenUrl="https://auth.vshn.net/token"
```


Command line flags will take precedence over both the configuration file and 
environment variables.


## Deploy

To deploy the `kube-token-refresher` you need OIDC credentials capable of 
requesting an access token, and Kubernetes credentials to `get` and `update`
the specified secret.

You can find an example deployment in `deploy/`.
