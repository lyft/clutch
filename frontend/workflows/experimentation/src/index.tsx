import type { WorkflowConfiguration } from "@clutch-sh/core";

import ListExperiments from "./list-experiments";
import { StartAbortExperiment, StartLatencyExperiment } from "./start-experiment";

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "experimentation",
    group: "Experimentation",
    displayName: "Experimentation",
    routes: {
      listExperiments: {
        path: "list",
        displayName: "List Experiments",
        description: "List Experiments.",
        component: ListExperiments,
      },
      startAbortExperiment: {
        path: "startabort",
        displayName: "Start an Abort Experiment",
        description: "Start an Abort Experiment.",
        component: StartAbortExperiment,
      },
      startLatencyExperiment: {
        path: "startlatency",
        displayName: "Start a Latency Experiment",
        description: "Start a Latency Experiment.",
        component: StartLatencyExperiment,
      },
    },
  };
};

export default register;
