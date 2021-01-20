import type { WorkflowConfiguration } from "@clutch-sh/core";

import { StartExperiment } from "./start-experiment";

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "server-experimentation",
    group: "Experimentation",
    displayName: "Server Fault Injection",
    routes: {
      startExperiment: {
        path: "start",
        displayName: "Start Experiment",
        description: "Start Server Experiment.",
        component: StartExperiment,
      },
    },
  };
};

export default register;
