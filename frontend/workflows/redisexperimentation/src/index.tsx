import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";

import { StartExperiment } from "./start-experiment";

export interface WorkflowProps extends BaseWorkflowProps {}

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "redis-experimentation",
    group: "Chaos Experimentation",
    displayName: "Redis Fault Injection",
    routes: {
      startExperiment: {
        path: "/start",
        displayName: "Start Experiment",
        description: "Start Redis Experiment.",
        component: StartExperiment,
      },
    },
  };
};

export default register;
