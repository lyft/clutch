import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";
import type { WizardChild } from "@clutch-sh/wizard";

import UpdateCapacity from "./update-capacity";

interface ResolverConfigProps {
    resolverType: string;
  }  

export interface WorkflowProps extends BaseWorkflowProps, ResolverConfigProps {}
export interface ResolverChild extends WizardChild, ResolverConfigProps {}

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "dynamodb",
    group: "AWS",
    displayName: "Dynamodb",
    routes: {
      updateCapacity: {
        path: "/capacity",
        component: UpdateCapacity,
        displayName: "Update Capacity",
        description: "Update the table or GSI provisioned capacity.",
        requiredConfigProps: ["resolverType", "notes"],
      },
    },
  };
};

export default register;