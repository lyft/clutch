import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";

export interface WorkflowProps extends BaseWorkflowProps {}

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "catalog",
    group: "Catalog",
    displayName: "Project Catalog",
    routes: {
      catalog: {
        path: "/",
        displayName: "Project Catalog",
        description: "Project Catalog",
        component: () => null,
        featureFlag: "projectCatalog",
      },
    },
  };
};

export default register;
