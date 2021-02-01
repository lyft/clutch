import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";

import Amiibo from "./hello-world";

export interface WorkflowProps extends BaseWorkflowProps {}

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@example.com",
    },
    path: "amiibo",
    group: "Amiibo",
    displayName: "Amiibo",
    routes: {
      landing: {
        path: "/lookup",
        description: "Lookup all Amiibo by name.",
        component: Amiibo,
      },
    },
  };
};

export default register;
