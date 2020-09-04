import type { WorkflowConfiguration } from "@clutch-sh/core";

import ListExperiments from "./list-experiments";
import ViewExperiment from "./view-experiment";

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
      viewExperiment: {
        path: "view/:id",
        displayName: "View Experiment",
        description: "View Experiment",
        component: ViewExperiment,
      }
    },
  };
};

export default register;
