import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";

import RemoteTriage from "./remote-triage";

export interface WorkflowProps extends BaseWorkflowProps {}

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "envoy",
    group: "Envoy",
    displayName: "Envoy",
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
