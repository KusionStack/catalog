import regex

schema SecretBase:
    """ SecretBase defines set of attributes shared by different secret types, e.g. SecretBasic
    and SecretOpaque. You can inherit this Schema to reuse these common attributes.

    Attributes
    ----------
    immutable: bool, default is Undefined, optional.
        Immutable, if set to true, ensures that data stored in the Secret cannot be updated.
    """

    # If immutable set to true, ensures that data stored in the Secret cannot be updated.
    immutable?:                 bool

schema BasicAuthSecret(SecretBase):
    """ BasicAuthSecret is used for username / password pairs, which commonly used for basic
    authentication scene.

    Attributes
    ----------
    type: str, default is basic, required.
        The actual type of secret.
    data: {str,str}, default is Undefined, optional.
        Data contains the username / password pairs. The key names must be username and password.

    Examples
    --------
    import models.schema.v1.workload.secret as sec

    basicAuth = sec.BasicAuthSecret {
        type: "basic"
        data: {
            "username": ""
            "password": ""
        }
    }
    """

    # The basic secret type.
    type:                       "basic"

    # Data contains the username / password pairs.
    # The key names must be username and password.
    data?:                      {str:str}

    check:
        all k in data {
            k in ["username", "password"]
        } if data, "a valid secret data key must in username and password for BasicAuthSecret"

schema TokenSecret(SecretBase):
    """ TokenSecret is used for generating a password or secure string used for passwords when
    the user is already known or not required.

    Attributes
    ----------
    type: str, default is token, required.
        The actual type of secret.
    params: {str,str}, default is Undefined, optional.
        Collection of parameters used to facilitate programmatic handling of secret data.
    data: {str,str}, default is Undefined, optional.
        Data contains token info if user already known. The key name must be token.

    Examples
    --------
    import models.schema.v1.workload.secret as sec

    token = sec.TokenSecret {
        type: "token"
        data: {
            # optional
            "token": ""
        }
    }
    """

    # The token secret type.
    type:                       "token"

    # Params defines extra parameters used to customize secret handling.
    params?:                    {str:str}

    # Data contains token info if user already known.
    data?:                      {str:str}

    check:
        all k in data {
            k in ["token"]
        } if data, "a valid secret data key name must be token for TokenSecret"

schema OpaqueSecret(SecretBase):
    """ OpaqueSecret is used for arbitrary key value pairs. This type of secrets are
    best used for allowing a user to input sensitive data at runtime.

    Attributes
    ----------
    type: str, default is opaque, required.
        The actual type of secret.
    data: {str,str}, default is Undefined, optional.
        Data contains arbitrary key value pairs.

    Examples
    --------
    import models.schema.v1.workload.secret as sec

    opaqueSec = sec.OpaqueSecret {
        type: "opaque"
    }
    """

    # The opaque secret type.
    type:                       "opaque"

    # Data contains arbitrary key value pairs.
    data?:                      {str:str}

    check:
        not data or all k in data {
            regex.match(k, r"[A-Za-z0-9_.-]*")
        }, "a valid secret data key must consist of alphanumeric characters, '-', '_' or '.'"

schema TLSSecret(SecretBase):
    """ TLSSecret is used for storing a certificate and its associated key that
    are typically used for TLS.

    Attributes
    ----------
    type: str, default is tls, required.
        The actual type of secret.
    data: {str,str}, default is Undefined, required.
        Data contains certificate and its associated key.

    Examples
    --------
    import models.schema.v1.workload.secret as sec

    tlsSecret = sec.TLSSecret {
        type: "tls"
        data: {
            # must be provided
            "tls.crt": ""
            "tls.key": ""
        }
    }
    """

    # The tls secret type.
    type:                       "tls"

    # Data contains certificate and its associated key.
    data:                       {str:str}

    check:
        all k in data {
            k in ["tls.crt", "tls.key"]
        } if data, "a valid secret data key name must be in tls.crt and tls.key for TLSSecret"
        # TODO: add check rule to ensure tls.crt and tls.key exist in data

schema ExternalSecret(SecretBase):
    """ ExternalSecret is used for retrieving secret data from external secret management systems, e.g.
    Hashicorp Vault, AWS Secrets Managers.

    Attributes
    ----------
    type: str, default is external, required.
        The actual type of secret.
    data: [ExternalSecretData], default is Undefined, required.
        Data defines the connection between the Kubernetes Secret keys and external secret provider data.

    Examples
    --------
    import models.schema.v1.workload.secret as sec

    externalSec = sec.ExternalSecret {
        type: "external"
        data: [
            sec.ExternalSecretData {
                secretKey: "username"
                name: "api-secret"
                property: "username"
            }
        ]
    }
    """

    # The external secret type.
    type:                       "external"

    # Data defines the connection between the Kubernetes Secret keys and external secret provider data.
    data:                       [ExternalSecretData]

schema ExternalSecretData:
    """ ExternalSecretData defines the connection between the Kubernetes Secret key (spec.data.<key>)
    and the external secret provider data.

    Attributes
    ----------
    secretKey: str, default is Undefined, required.
        SecretKey defines the key in Kubernetes secret data.
    name: str, default is Undefined, required.
        Name specifies the name of the secret in external secret provider to read.
    property: str, default is Undefined, optional.
        Used to select a specific property of the secret data (if a map),
        if target external secret provider supported.
    version: str, default is Undefined, optional.
        Specifies the version of the secret to return, if target external secret provider supported.

    Examples
    --------
    import models.schema.v1.workload.secret as sec

    esData = sec.ExternalSecretData {
        secretKey: "username"
        name: "api-secret"
        property: "username"
    }
    """

    # SecretKey defines the key in Kubernetes secret data.
    secretKey:                  str

    # Name specifies the name of the secret in external secret provider to read.
    name:                       str

    # Used to select a specific property of the secret data (if a map), if supported.
    property?:                  str

    # Specifies the version of the secret to return, if supported.
    version?:                   str

    check:
        regex.match(secretKey, r"[A-Za-z0-9_.-]*"), "a valid secret data key must consist of alphanumeric characters, '-', '_' or '.'"
