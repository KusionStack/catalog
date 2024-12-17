schema TypedLocalObjectReference:
    """ TypedLocalObjectReference contains enough information to let you locate the typed referenced object inside the
    same namespace.

    Attributes
    ----------
    kind: str, default is Undefined, required.
        Kind is the type of resource being referenced.
    name: str, default is Undefined, required.
        Name is the name of resource being referenced.
    apiGroup: str, optional.
        APIGroup is the group for the resource being referenced. If APIGroup is not specified, the specified Kind must
        be in the core API group. For any other third-party types, APIGroup is required.
    """

    # Kind is the type of resource being referenced.
    kind:                      str

    # Name is the name of resource being referenced.
    name:                      str

    # APIGroup is the group for the resource being referenced. If APIGroup is not specified, the specified Kind must
    # be in the core API group. For any other third-party types, APIGroup is required.
    apiGroup?:                 str