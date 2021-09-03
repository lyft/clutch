import type { BaseWorkflowProps, DefaultWorkflowConfig } from "@clutch-sh/core";

import Card from "./card";
import Dash from "./dash";
import { useDashState } from "./dash-hooks";

export interface WorkflowProps extends BaseWorkflowProps {}

const register = (): DefaultWorkflowConfig => {
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

export type { DashState } from "./types";
export { Card, Dash, useDashState };
