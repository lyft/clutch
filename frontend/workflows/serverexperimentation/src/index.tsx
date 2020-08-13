import type { WorkflowConfiguration } from "@clutch-sh/core";

import { StartAbortExperiment, StartLatencyExperiment } from "./start-experiment";

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "serverexperimentation",
    group: "Experimentation",
    displayName: "Server Experimentation",
    routes: {
      startAbortExperiment: {
        path: "startabort",
        displayName: "Start Abort Experiment",
        description: "Start Abort Experiment.",
        component: StartAbortExperiment,
      },
      startLatencyExperiment: {
        path: "startlatency",
        displayName: "Start Latency Experiment",
        description: "Start Latency Experiment.",
        component: StartLatencyExperiment,
      },
    },
  };
};

export default register;
