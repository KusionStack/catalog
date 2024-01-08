schema MySQL: 
    """ MySQL describes the attributes to locally deploy or create a cloud provider
    managed mysql database instance for the workload. 

    Attributes
    ----------
    type: "local" | "cloud", defaults to Undefined, required. 
        Type defines whether the mysql database is deployed locally or provided by 
        cloud vendor. 
    version: str, defaults to Undefined, required. 
        Version defines the mysql version to use. 

    Examples
    --------
    Instantiate a local mysql database with version of 5.7. 

    import models.schema.v1.accessories.mysql

    mysql: mysql.MySQL {
        type: "local"
        version: "5.7"
    }
    """

    # The deployment mode of the mysql database. 
    type:       "local" | "cloud"

    # The mysql database version to use. 
    version:    str
    