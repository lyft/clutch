import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";

import Catalog from "./catalog";
import Details from "./details";

export interface WorkflowProps extends BaseWorkflowProps {}

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "catalog",
    group: "Catalog",
    displayName: "Catalog",
    routes: {
      catalog: {
        path: "/",
        displayName: "Project Catalog",
        description: "A searchable catalog of services",
        component: Catalog,
        // featureFlag: "projectCatalog",
      },
      details: {
        path: "/:projectId",
        description: "Service Detail View",
        component: Details,
        // featureFlag: "projectCatalog",
      },
    },
  };
};

export default register;
