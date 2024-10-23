import type { BaseWorkflowProps, NoteConfig, WorkflowConfiguration } from "@clutch-sh/core";
import type { WizardChild } from "@clutch-sh/wizard";

import UpdateCapacity from "./update-capacity";

interface ResolverConfigProps {
  resolverType: string;
  notes?: NoteConfig[];
}

interface TableDetailsProps {
  enableOverride?: boolean;
  notes?: NoteConfig[];
}

export interface WorkflowProps extends BaseWorkflowProps, ResolverConfigProps, TableDetailsProps {}
export interface ResolverChild extends WizardChild, ResolverConfigProps {}
export interface TableDetailsChild extends WizardChild, TableDetailsProps {}

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:clutch@lyft.com",
    },
    path: "dynamodb",
    group: "AWS",
    displayName: "DynamoDB",
    routes: {
      updateCapacity: {
        path: "/capacity",
        component: UpdateCapacity,
        displayName: "Update Capacity",
        description: "Update the table or GSI provisioned capacity.",
        requiredConfigProps: ["resolverType"],
        layoutProps: {
          variant: "wizard",
        },
      },
    },
  };
};

export default register;
