import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";
import { Dash, useDashState } from "./dash";

import { Group } from "./types";

export interface WorkflowProps extends BaseWorkflowProps {}

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "clutch@lyft.com",
      contactUrl: "mailto:clutch@lyft.com",
    },
    path: "dash",
    group: "Dash",
    displayName: "Dash",
    routes: {
      landing: {
        path: "/",
        displayName: "Dash",
        description: "Display helpful information for multiple services.",
        component: Dash,
      },
    },
  };
};

export default register;

export { Dash, useDashState };
