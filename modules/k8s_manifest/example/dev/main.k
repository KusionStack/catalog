import kam.v1.app_configuration as ac
import k8s_manifest

test: ac.AppConfiguration {
    accessories: {
        "k8s_manifests": k8s_manifest.K8sManifest {
            paths: [
                "./test.yaml"
            ]
        }
    }
}
