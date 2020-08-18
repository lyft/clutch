import type { WorkflowConfiguration } from "@clutch-sh/core";

import ListExperiments from "./list-experiments";

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
    },
  };
};

export default register;
