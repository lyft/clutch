Kubernetes Example
===

If you dont have a Kubernetes cluster avaiable you can use [minikube](https://minikube.sigs.k8s.io/docs/start/)
to provision a local cluster.

To deploy a working clutch deployment on an existing Kubernetes cluster:

1. Install kubectl and set up `~/.kube/config` to point to your cluster.
1. Run the following command to create a clutch namespace and all the necessary clutch components

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/lyft/clutch/main/examples/kubernetes/clutch.yaml
    
    # OR
    # if you are working off a local copy of this repo
    kubectl apply -f examples/kubernetes/clutch.yaml
    ```
1. Connect to your new clutch deployment. You can follow your cloud provider's documentation for configuring ingress on your cluster, or use port-forwarding: 
    ```
    kubectl port-forward service/clutch 8080:8080 -n clutch
    ```
1. The clutch UI should be available in your browser at http://localhost:8080
