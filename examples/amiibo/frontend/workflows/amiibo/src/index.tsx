import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch/core";

import HelloWorld from "./hello-world";

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
        path: "/",
        description: "Lookup all Amiibo by name.",
        component: HelloWorld,
      },
    },
  };
};

export default register;
