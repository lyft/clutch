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
    displayName: "Experimentation",
    routes: {
      listExperiments: {
        path: "list",
        displayName: "List Experiments",
        description: "List Experiments.",
        component: ListExperiments,
      },
      viewExperimentRun: {
        path: "view/:id",
        displayName: "View Experiment Run",
        description: "View Experiment Run",
        hideNav: true,
        component: ViewExperimentRun,
      },
    },
  };
};

export default register;
