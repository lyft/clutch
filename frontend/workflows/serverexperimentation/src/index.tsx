import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";
import { StartAbortExperiment, StartLatencyExperiment } from "./start-experiment";

export interface WorkflowProps extends BaseWorkflowProps {}

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "hello@lyft.com",
      contactUrl: "mailto:example@lyft.com",
    },
    path: "serverexperimentation",
    group: "Experimentation Server",
    displayName: "Experimentation Server",
    routes: {
      startAbortExperiment: {
        path: "startabort",
        description: "Start Abort Experiment.",
        displayName: "Start Abort Experiment",
        component: StartAbortExperiment,
      },
      startLatencyExperiment: {
        path: "startlatency",
        description: "Start Latency Experiment.",
        displayName: "Start Latency Experiment",
        component: StartLatencyExperiment,
      }
    },
  };
};

export default register;

