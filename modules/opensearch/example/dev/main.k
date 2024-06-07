# The configuration codes in perspective of developers. 
import kam.v1.app_configuration as ac
import service as svc
import service.secret as sec
import service.container as c
import opensearch.OpenSearch as o

agent: ac.AppConfiguration {
    # Declare the workload configurations.
    workload: svc.Service {
        containers: {
            opensearch: c.Container {
                image: "kusionstack/ai-agent-demo:v0.1.0"
            }
        }
    }
    # Declare the openSearch module configurations.
    accessories: {
        "opensearch": o.OpenSearch {
            domainName: "agent-test"
            engineVersion: "OpenSearch_2.13"
        }
    }
}
