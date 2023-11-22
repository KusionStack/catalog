import regex

schema Secret:
    """ Secrets are used to provide data that is considered sensitive like passwords, API keys,
    TLS certificates, tokens or other credentials.

    Attributes
    ----------
    type: str, default is Undefined, required.
        Type of secret, used to facilitate programmatic handling of secret data.
    data: {str,str}, default is Undefined, optional.
        Data contains the non-binary secret data in string form.
    immutable: bool, default is Undefined, optional.
        Immutable, if set to true, ensures that data stored in the Secret cannot be updated.

    Examples
    --------
    import models.schema.v1.workload.secret as sec

    basicAuth = sec.Secret {
        type: "basic"
        data: {
            "username": ""
            "password": ""
        }
    }
    """

    # The type of target secret.
    type:                       "basic" | "token" | "opaque" | "certificate"

    # Data defines the keys and data that will be used by secret.
    data?:                      {str:str}

    # If immutable set to true, ensures that data stored in the Secret cannot be updated.
    immutable?:                 bool

    check:
        not data or all k in data {
            regex.match(k, r"[A-Za-z0-9_.-]*")
        }, "a valid secret data key must consist of alphanumeric characters, '-', '_' or '.'"
