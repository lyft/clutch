import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";
import type { WizardChild } from "@clutch-sh/wizard";

import UpdateShardCount from "./update-shard-count";

interface ConfigurationProps {
  resolverType: string;
}

interface WizardFeedbackProps {
  enableFeedback?: boolean;
}

export interface WorkflowProps extends BaseWorkflowProps, ConfigurationProps, WizardFeedbackProps {}
export interface ResolverChild extends WizardChild, ConfigurationProps {}

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "kinesis",
    group: "AWS",
    displayName: "Kinesis",
    routes: {
      updateShardCount: {
        path: "shards",
        component: UpdateShardCount,
        displayName: "Shard Count",
        description: "Update Kinesis stream shard count.",
        requiredConfigProps: ["resolverType"],
      },
    },
  };
};

export default register;
