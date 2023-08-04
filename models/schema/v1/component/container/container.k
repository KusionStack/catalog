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

    check:
        all e in env {
            regex.match(e, r"^[-._a-zA-Z][-._a-zA-Z0-9]*$")
        } if env, "a valid environment variable name must consist of alphabetic characters, digits, '_', '-', or '.', and must not start with a digit"
