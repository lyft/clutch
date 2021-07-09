import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";
import { Dash } from "./dash";

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
        component: Dash,
      },
    },
  };
};

export default register;
