import type { WorkflowConfiguration } from "@clutch-sh/core";

import Config from "./details/config";
import Catalog from "./catalog";
import Details from "./details";

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
        layoutProps: {
          variant: "standard",
          title: "Project Catalog",
          subtitle: "A catalog of all projects.",
        },
      },
      details: {
        path: "/:projectId",
        description: "Service Detail View",
        component: Details,
        featureFlag: "projectCatalog",
        layoutProps: {
          variant: "standard",
          onlyBreadcrumbs: true,
        },
      },
      configLanding: {
        path: "/:projectId/config",
        description: "Project Configuration Landing",
        component: Config,
        featureFlag: "projectCatalog",
        layoutProps: {
          variant: "standard",
          onlyBreadcrumbs: true,
        },
      },
      configPage: {
        path: "/:projectId/config/:configType",
        description: "Project Configuration Page",
        component: Config,
        featureFlag: "projectCatalog",
        layoutProps: {
          variant: "standard",
          onlyBreadcrumbs: true,
        },
      },
    },
  };
};

export * from "./helpers";
export { CardType, DynamicCard, MetaCard } from "./details/components/card";
export { useProjectDetailsContext } from "./details/context";
export { Details as ProjectDetails, Config as ProjectConfig };
export type { CatalogDetailsCard } from "./details/components/card";
export type { ProjectInfoChip } from "./details/info/chipsRow";
export type { CatalogDetailsChild, ProjectConfigLink, DetailsLayoutOptions } from "./types";
export type { ProjectConfigProps } from "./details/config";

export default register;
