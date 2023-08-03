import models.schema.v1.component.container as c

schema Component:
    """ Component defines the delivery artifact of one application. Each component
    can be composed of a Workload and its supporting Accessories.

    Attributes
    ----------
    containers: {str:c.Container}, default is Undefined, required.
        Containers defines the templates of containers to be ran.
        More info: https://kubernetes.io/docs/concepts/containers
    replicas: int, default is 2, required.
        Number of container replicas that should be ran.
    labels: {str:str}, default is Undefined, optional.
        Labels are key/value pairs that are attached to the component.
    annotations: {str:str}, default is Undefined, optional.
        Annotations are key/value pairs that attach arbitrary non-identifying metadata to the component.

    Examples
    --------
    Instantiate a long-running service Component with "nginx:v1" image and customized labels.

    component = cmp.Component {
        containers: {
            "nginx": c.Container {
                image: "nginx:v1"
            }
        }
        labels: {
            "key": "value"
        }
    }
    """

    # The templates of containers to be ran.
    containers:                 {str:c.Container}

    # The number of containers that should be ran.
    # Default is 2 to meet high availability requirements.
    # Not supported when type is Job.
    replicas:                   int = 2

    ###### Other metadata info
    # Labels and annotations can be used to attach arbitrary metadata as key-value pairs to resources.
    labels?:                    {str:str}
    annotations?:               {str:str}
