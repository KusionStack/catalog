# The configuration codes in perspective of developers. 
import kam.v1.app_configuration as ac
import service
import service.container as c
import inference.v1.inference
import network as n

inference: ac.AppConfiguration {
    # Declare the workload configurations. 
    workload: service.Service {
        containers: {
            myct: c.Container {image: "kangy126/app"}
        }
        replicas: 1
    }
    # Declare the inference module configurations. 
    accessories: {
        "inference": inference.Inference {
            model: "llama3"
            framework: "Ollama"
        }
        "network": n.Network {ports: [n.Port {
            port: 80
            targetPort: 5000
            public: True
        }]}
    }
}