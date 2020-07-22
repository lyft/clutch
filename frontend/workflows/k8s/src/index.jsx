import DeletePod from "./delete-pod";
import ResizeHPA from "./resize-hpa";

const register = () => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "k8s",
    group: "K8s",
    displayName: "K8s",
    routes: {
      deletePod: {
        path: "pod/delete",
        displayName: "Delete Pod",
        description: "Delete a K8s pod.",
        component: DeletePod,
        requiredConfigProps: ["resolverType"],
      },
      resizeHPA: {
        path: "hpa/resize",
        displayName: "Resize HPA",
        description: "Resize a horizontal autoscaler.",
        component: ResizeHPA,
        requiredConfigProps: ["resolverType"],
      },
    },
  };
};

export default register;
