import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";
import type { WizardChild } from "@clutch-sh/wizard";

import ChangeCapacity from "./change-capacity";

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
      changeCapacity: {
        path: "capacity",
        component: ChangeCapacity,
        displayName: "Change Capacity",
        description: "Change the table or GSI provisioned capacity.",
        requiredConfigProps: ["resolverType"],
      },
    },
  };
};

export default register;
