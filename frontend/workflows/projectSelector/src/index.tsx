import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";

import HelloWorld from "./hello-world";

export interface WorkflowProps extends BaseWorkflowProps {}

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "hello@example.com",
      contactUrl: "mailto:hello@example.com",
    },
    path: "projectselector",
    group: "Project Selector",
    displayName: "Project Selector",
    routes: {
      landing: {
        path: "/",
        description: "Filter your projects.",
        component: HelloWorld,
      },
    },
  };
};

export default register;
