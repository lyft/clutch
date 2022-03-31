import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";

import type { CatalogDetailsCard } from "./details/card";
import { DynamicCard, MetaCard } from "./details/card";
import type { ProjectInfoChip } from "./details/info/chipsRow";
import Catalog from "./catalog";
import Details from "./details";

type DetailCard = CatalogDetailsCard | typeof DynamicCard | typeof MetaCard;

type CatalogDetailsChild = React.ReactElement<DetailCard>;

export interface WorkflowProps extends BaseWorkflowProps {}
export interface ProjectDetailsWorkflowProps {
  children?: CatalogDetailsChild | CatalogDetailsChild[];
  chips?: ProjectInfoChip[];
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
    },
  };
};

export { DynamicCard, MetaCard };
export { LastEvent } from "./details/helpers";
export { useProjectDetailsContext } from "./details/context";
export { Details as ProjectDetails };
export type { CatalogDetailsCard, CatalogDetailsChild, ProjectInfoChip };

export default register;
