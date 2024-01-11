schema Prometheus:
    """ Prometheus can be used to define monitoring requirements

    Attributes
    ----------
    path: str, default is /metrics, optional
        The path to scrape metrics from.
    port: str, default is container ports when scraping pod (monitorType is pod) and service port when scraping service (monitorType is service), optional
        The port to scrape metrics from. When using Prometheus operator, this needs to be the port NAME. Otherwise, this can be a port name or a number.

    Examples
    --------
    import catalog.models.schema.v1.monitoring as m

    monitoring: m.Prometheus{
        path:           "/metrics"
        port:           "web"
    }
    """

    # Path defines the path from which Prometheus scrapes the target.
    path?:                      str

    # Port defines the port from which Prometheus scrapes the target.
    port?:                      str