schema Volume:
    """ On-disk files in container are ephemeral, Volume can be used to store persistent data
    that can be mounted by containers.

    Attributes
    ----------
    size: str | int, default is opaque, required.
        Size configures the default capacity of the target volume.
    accessModes: [str], default is Undefined, optional.
        AccessModes contains the desired access modes the volume should have.
        More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes
    class: str, default is Undefined, optional.
        Name of StorageClass required by the volume.

    Examples
    --------
    import models.schema.v1.workload.volume as v

    vol = v.Volume {
        size: 100
        accessModes: ["ReadWriteOnce"]
    }
    """

    # Size configures the default capacity of the target volume.
    size:                       str | int

    # AccessModes contains the desired access modes the volume should have.
    accessModes:                [str]

    # Name of StorageClass required by the volume.
    class?:                     str

    check:
        len(class) < 253 if class, "class name must be no more than 253 characters"
