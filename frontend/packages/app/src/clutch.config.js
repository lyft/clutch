module.exports = {
  "@clutch-sh/ec2": {
    terminateInstance: {
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
    rebootInstance: {
      hideNav: true,
      componentProps: {
        resolverType: "clutch.aws.ec2.v1.Instance",
      },
    },
    resizeAutoscalingGroup: {
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
  },
  "@clutch-sh/envoy": {
    remoteTriage: {
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
  },
  "@clutch-sh/k8s": {
    deletePod: {
      componentProps: {
        resolverType: "clutch.k8s.v1.Pod",
      },
    },
    resizeHPA: {
      trending: true,
      componentProps: {
        resolverType: "clutch.k8s.v1.HPA",
      },
    },
    kubeDashboard: {
      trending: true,
    },
    scaleResources: {
      trending: true,
      componentProps: {
        resolverType: "clutch.k8s.v1.Deployment",
      },
    },
    cordonNode: {
      trending: true,
      componentProps: {
        resolverType: "clutch.k8s.v1.Node",
      },
    },
  },
  "@clutch-sh/project-catalog": {
    catalog: {
      trending: false,
      hideNav: true,
    },
    details: {
      hideNav: true,
    },
    configLanding: {
      hideNav: true,
      componentProps: {
        defaultRoute: "/",
      },
    },
    configPage: {
      hideNav: true,
      componentProps: {
        defaultRoute: "/",
      },
    },
  },
};
