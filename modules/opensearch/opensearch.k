schema OpenSearch:
    """ OpenSearch is a module schema of OpenSearch. Currently, it only supports AWS OpenSearch Service

    Attributes
    ----------
    engineVersion: str, default is Undefined, optional. 
            Either Elasticsearch_X.Y or OpenSearch_X.Y to specify the engine version for the Amazon OpenSearch Service domain. For example, OpenSearch_1.0 or Elasticsearch_7.9. Defaults to the latest version of OpenSearch.
    domainName: str, default is Undefined, required. 
        Name of the domain. 
    
    Examples
    --------
    import opensearch.OpenSearch as o

    accessories: {
        "opensearch": o.OpenSearch {
            domainName: "example"
            engineVersion: "OpenSearch_1.0"
        }
    }
    """

    # Name of the domain
    domainName :             str

    # Either Elasticsearch_X.Y or OpenSearch_X.Y to specify the engine version for the Amazon OpenSearch Service domain. For example, OpenSearch_1.0 or Elasticsearch_7.9. Defaults to the lastest version of OpenSearch.
    engineVersion?:          str