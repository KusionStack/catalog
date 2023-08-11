import regex

schema Container:
    """ Container describes how the Application's tasks are expected to be run. Depending on
    the replicas parameter 1 or more containers can be created from each template.

    Attributes
    ----------
    image: str, default is Undefined, required.
        Image refers to the Docker image name to run for this container.
        More info: https://kubernetes.io/docs/concepts/containers/images
    command: [str], default is Undefined, optional.
        Entrypoint array. Not executed within a shell.
        Command will overwrite the ENTRYPOINT value set in the Dockfile, otherwise the Docker
        image's ENTRYPOINT is used if this is not provided.
    args: [str], default is Undefined, optional.
        Arguments to the entrypoint.
        Args will overwrite the CMD value set in the Dockfile, otherwise the Docker
        image's CMD is used if this is not provided.
    env: {str:str}, default is Undefined, optional.
        List of environment variables to set in the container.
        The value of the environment variable may be static text or a value from a secret.
    workingDir: str, default is Undefined, optional.
        The working directory of the running process defined in entrypoint.
        Continer runtime's default will be used if this is not specified.

     Examples
    --------
    import models.schema.v1.component.container as c

    web = c.Container {
        image:   "nginx:latest"
        command: ["/bin/sh", "-c", "echo hi"]
        env: {
            "name": "value"
        }
    }
    """

    # Image to run for this container
    image:                      str

    # Entrypoint array.
    # The image's ENTRYPOINT is used if this is not provided.
    command?:                   [str]

    # Arguments to the entrypoint.
    # The image's CMD is used if this is not provided.
    args?:                      [str]

    # Collection of environment variables to set in the container.
    # The value of environment variable may be static text or a value from a secret.
    env?:                       {str:str}

    # The current working directory of the running process defined in entrypoint.
    workingDir?:                str

    # Dirs configures one or more volumes to be mounted to the specified folder.
    dirs?:                      {str:str}
    # Files configures one or more files to be created in the container.
    files?:                     {str:FileSpec}

    check:
        all e in env {
            regex.match(e, r"^[-._a-zA-Z][-._a-zA-Z0-9]*$")
        } if env, "a valid environment variable name must consist of alphabetic characters, digits, '_', '-', or '.', and must not start with a digit"

schema FileSpec:
    """ FileSpec defines the target file in a Container.

    Attributes
    ----------
    content: str, default is Undefined, optional.
        File content in plain text.
    contentFrom: str, default is Undefined, optional.
        Source for the file content, reference to a secret of configmap value.
    mode: str, default is Undefined, optional.
        Mode bits used to set permissions on this file, must be an octal value
        between 0000 and 0777 or a decimal value between 0 and 511

    Examples
    --------
    import models.schema.v1.component.container as c

    tmpFile = c.FileSpec {
        content: "some file contents"
        mode: "0777"
    }
    """

    # The content of target file in plain text.
    content?:                   str

    # Source for the file content, might be a reference to a secret value.
    contentFrom?:               str

    # Mode bits used to set permissions on this file.
    # Defaults to 0644.
    mode?:                      str = "0644"

    check:
        not content or not contentFrom, "content and contentFrom are mutually exclusive"
        regex.match(mode, r"^[0-7]{3,4}$"), "valid mode must between 0000 and 0777, both inclusive"
