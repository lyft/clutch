import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";
import type { WizardChild } from "@clutch-sh/wizard";

import RemoteTriage from "./remote-triage";

interface TriageProps {
  host?: string;
}

export interface WorkflowProps extends BaseWorkflowProps {}
export interface TriageChild extends WizardChild, TriageProps {}

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "envoy",
    group: "Envoy",
    displayName: "Envoy",
    defaultLayoutProps: {
      variant: "wizard",
    },
    routes: {
      remoteTriage: {
        path: "triage",
        component: RemoteTriage,
        displayName: "Remote Triage",
        description: "Triage Envoy configurations.",
      },
    },
  };
};

export default register;
