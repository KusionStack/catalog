import models.schema.v1.component as c

schema AppConfiguration:
    """ AppConfiguration is a developer-centric definition that describes how to run an Application.
    This application model builds upon a decade of experience at AntGroup running super large scale
    internal developer platform, combined with best-of-breed ideas and practices from the community.

    Attributes
    ----------
    components: {str:c.Component}, default is Undefined, required.
        Component defines the delivery artifact of one application. Each application can be
        composed by multiple components.

    Examples
    --------
    Instantiate an App with a long-running service and its image is "nginx:v1"

    appConfiguration = ac.AppConfiguration {
        components: {
            "proxy": cmp.Component {
                containers: {
                    "nginx": c.Container {
                        image: "nginx:v1"
                    }
                }
            }
        }
    }
    """

    # Component defines the delivery artifact of one application.
    # Each application can be composed by multiple components.
    components?:                 {str:c.Component}
