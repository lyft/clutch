module.exports = {
    "@clutch-sh/ec2": {
        terminateInstance: {
            trending: true,
            componentProps: {
                resolverType: "clutch.aws.ec2.v1.Instance",
            },
        },
        resizeAutoscalingGroup: {
            trending: true,
            componentProps: {
                resolverType: "clutch.aws.ec2.v1.AutoscalingGroup",
            },
        },
    },
    "@clutch-sh/envoy": {
        remoteTriage: {
            trending: true,
            componentProps: {
                options: {
                    "Clusters": "clusters",
                    "Listeners": "listeners",
                    "Runtime": "runtime",
                    "Stats": "stats",
                    "Server Info": "serverInfo",
                },
            },
        },
    },
    "@{{ .RepoName }}/echo": {
        echo: {
            trending: true,
        }
    }
};
