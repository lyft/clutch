const contextValues: { workflows: any[] } = {
  workflows: [
    {
      developer: {
        name: "Lyft",
        contactUrl: "mailto:hello@clutch.sh",
      },
      path: "ec2",
      group: "AWS",
      displayName: "EC2",
      routes: [
        {
          path: "instance/terminate",
          displayName: "Terminate Instance",
          description: "Terminate an EC2 instance.",
          requiredConfigProps: ["resolverType"],
          trending: true,
          componentProps: {
            resolverType: "clutch.aws.ec2.v1.Instance",
            notes: [
              {
                severity: "info",
                text: "The instance may take several minutes to shut down.",
              },
            ],
          },
        },
        {
          path: "asg/resize",
          displayName: "Resize Autoscaling Group",
          description: "Resize an autoscaling group.",
          requiredConfigProps: ["resolverType"],
          componentProps: {
            resolverType: "clutch.aws.ec2.v1.AutoscalingGroup",
            notes: [
              {
                severity: "info",
                text:
                  "The autoscaling group may take several minutes to bring additional instances online.",
              },
            ],
          },
        },
      ],
    },
    {
      developer: {
        name: "Lyft",
        contactUrl: "mailto:hello@clutch.sh",
      },
      path: "envoy",
      group: "Envoy",
      displayName: "Envoy",
      routes: [
        {
          path: "triage",
          displayName: "Remote Triage",
          description: "Triage Envoy configurations.",
          trending: true,
          componentProps: {
            options: {
              Clusters: "clusters",
              Listeners: "listeners",
              Runtime: "runtime",
              Stats: "stats",
              "Server Info": "serverInfo",
            },
          },
        },
      ],
    },
    {
      developer: {
        name: "Lyft",
        contactUrl: "mailto:hello@clutch.sh",
      },
      path: "k8s",
      group: "K8s",
      displayName: "K8s",
      routes: [
        {
          path: "pod/delete",
          displayName: "Delete Pod",
          description: "Delete a K8s pod.",
          requiredConfigProps: ["resolverType"],
          componentProps: {
            resolverType: "clutch.k8s.v1.Pod",
          },
        },
        {
          path: "hpa/resize",
          displayName: "Resize HPA",
          description: "Resize a horizontal autoscaler.",
          requiredConfigProps: ["resolverType"],
          trending: true,
          componentProps: {
            resolverType: "clutch.k8s.v1.HPA",
          },
        },
        {
          path: "dashboard",
          displayName: "Kubernetes Dashboard",
          description: "Dashboard for Kubernetes Resources.",
          requiredConfigProps: [],
          trending: true,
        },
        {
          path: "node/cordon",
          displayName: "Cordon/Uncordon Node",
          description: "Cordon or uncordon a node",
          requiredConfigProps: ["resolverType"],
          trending: true,
          componentProps: {
            resolverType: "clutch.k8s.v1.Node",
          },
        },
        {
          path: "probe/update",
          displayName: "Update Probes",
          description: "Update Probes on deployments",
          requiredConfigProps: ["resolverType"],
          trending: true,
          componentProps: {
            resolverType: "clutch.k8s.v1.Deployment",
          },
        },
      ],
    },
  ],
};

export default contextValues;
