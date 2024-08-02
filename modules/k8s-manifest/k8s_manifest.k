schema K8sManifest: 
    """ K8sManifest defines the paths of the YAML files, or the directories of the raw Kubernetes
    manifests, which will be jointly appended to the Resources of Spec. 

    Attributes
    ----------
    paths: [str], default is Undefined, required. 
        The paths of the YAML files, or the directories of the raw Kubernetes manifests. 
    
    Examples
    --------
    import k8s_manifest

    accessories: {
        "k8s_manifest": k8s_manifest.K8sManifest {
            paths: [
                # The path of a YAML file. 
                "/path/to/my/k8s-manifest.yaml", 
                # The path of a directory containing K8s manifests. 
                "/dir/to/my/k8s-manifests"
            ]
        }
    }
    """

    # The paths of the YAML files, or the directories of the raw Kubernetes manifests. 
    paths: [str] 

    check:
        len(paths) > 0,     "paths must be specified"
