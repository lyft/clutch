import type { WorkflowConfiguration } from "@clutch-sh/core";

import FormFields from "./core/form-fields";
import PageLayout from "./core/page-layout";
import ListExperiments from "./list-experiments";
import ViewExperimentRun from "./view-experiment-run";

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "experimentation",
    group: "Chaos Experimentation",
    displayName: "Chaos Experimentation",
    routes: {
      listExperiments: {
        path: "list",
        displayName: "Manage Experiments",
        description: "Manage experiments.",
        component: ListExperiments,
        layoutProps: {
          variant: "standard",
        },
      },
      viewExperimentRun: {
        path: "run/:runID",
        displayName: "View Experiment Run",
        description: "View Experiment Run.",
        hideNav: true,
        component: ViewExperimentRun,
        layoutProps: {
          variant: "standard",
        },
      },
    },
  };
};

export default register;
export type { FormItem } from "./core/form-fields";
export { PageLayout, FormFields };
