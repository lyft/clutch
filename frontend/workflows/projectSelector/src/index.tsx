import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";

import ProjectSelector from "./project-selector";

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
        displayName: "Project Selector",
        description: "Filter your projects.",
        component: ProjectSelector,
      },
    },
  };
};

export default register;

export { ProjectSelector };