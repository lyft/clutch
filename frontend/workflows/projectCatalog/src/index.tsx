import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";

import Catalog from "./catalog";

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
        featureFlag: "projectCatalog",
      },
      details: {
        path: "/:project",
        description: "Service Detail View",
        component: () => null,
        featureFlag: "projectCatalog",
      },
    },
  };
};

export default register;
