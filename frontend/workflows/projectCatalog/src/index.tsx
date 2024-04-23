import type { BaseWorkflowProps, ClutchError, WorkflowConfiguration } from "@clutch-sh/core";

import type { CatalogDetailsCard } from "./details/card";
import { CardType, DynamicCard, MetaCard } from "./details/card";
import type { ProjectInfoChip } from "./details/info/chipsRow";
import Catalog from "./catalog";
import Config from "./config";
import Details from "./details";

type DetailCard = CatalogDetailsCard | typeof DynamicCard | typeof MetaCard;

interface ProjectCatalogProps {
  allowDisabled?: boolean;
}

export interface ProjectConfigLink {
  title: string;
  path: string;
  icon?: React.ReactElement;
}

export interface ProjectConfigProps {
  title: string;
  path: string;
  onError?: (error: ClutchError) => void;
}

type CatalogDetailsChild = React.ReactElement<DetailCard>;

export type ProjectConfigPage = React.ReactElement<ProjectConfigProps>;

export interface WorkflowProps extends BaseWorkflowProps, ProjectCatalogProps {}

export interface ProjectDetailsWorkflowProps extends WorkflowProps, ProjectCatalogProps {
  children?: CatalogDetailsChild | CatalogDetailsChild[];
  chips?: ProjectInfoChip[];
  configLinks?: ProjectConfigLink[];
}

export interface ProjectDetailsConfigWorkflowProps extends WorkflowProps, ProjectCatalogProps {
  children?: ProjectConfigPage | ProjectConfigPage[];
  defaultRoute?: string;
}

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

export { CardType, DynamicCard, MetaCard };
export { LastEvent } from "./details/helpers";
export { useProjectDetailsContext } from "./details/context";
export { Details as ProjectDetails, Config as ProjectConfig };
export type { CatalogDetailsCard, CatalogDetailsChild, ProjectInfoChip };

export default register;
