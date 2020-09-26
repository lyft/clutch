import type { WorkflowConfiguration } from "@clutch-sh/core";

import ListExperiments from "./list-experiments";
import ViewExperimentRun from "./view-experiment-run";

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "experimentation",
    group: "Experimentation",
    displayName: "Fault Injection",
    routes: {
      listExperiments: {
        path: "list",
        displayName: "Experiments",
        description: "Manage fault injection experiments.",
        component: ListExperiments,
      },
      viewExperimentRun: {
        path: "run/:runID",
        displayName: "View Experiment Run",
        description: "View Experiment Run.",
        hideNav: true,
        component: ViewExperimentRun,
      },
    },
  };
};

export default register;
