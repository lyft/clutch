module.exports = {
    "clutch-ec2": {
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
    "clutch-envoy": {
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
    "clutch-example": {
        echo: {
            trending: true,
        }
    }
};
