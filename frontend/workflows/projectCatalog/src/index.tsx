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
      },
      details: {
        path: "/:projectId",
        description: "Service Detail View",
        component: Details,
        featureFlag: "projectCatalog",
      },
      configLanding: {
        path: "/:projectId/config",
        description: "Project Configuration Landing",
        component: Config,
        featureFlag: "projectCatalog",
      },
      configPage: {
        path: "/:projectId/config/:configType",
        description: "Project Configuration Page",
        component: Config,
        featureFlag: "projectCatalog",
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
export type { CatalogDetailsChild } from "./types";

export default register;
