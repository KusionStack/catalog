schema PostgreSQL: 
    """ PostgreSQL describes the attributes to locally deploy or create a cloud provider
    managed postgresql database instance for the workload. 

    Attributes
    ----------
    type: "local" | "cloud", defaults to Undefined, required. 
        Type defines whether the postgresql database is deployed locally or provided by
        cloud vendor. 
    version: str, defaults to Undefined, required. 
        Version defines the postgres version to use. 

    Examples
    --------
    Instantiate a local postgresql database with image version of 14.0. 

    import catalog.models.schema.v1.accessories.postgres
    
    accessories: {
        "postgres": postgres.PostgreSQL {
            type:   "local"
            version: "14.0"
        }
    }
    """

    # The deployment mode of the postgresql database. 
    type:       "local" | "cloud"

    # The postgresql database version to use. 
    version:    str
    