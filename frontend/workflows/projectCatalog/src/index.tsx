import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";

import type { DetailsCard } from "./details/card";
import type { ProjectInfoChip } from "./details/info/types";
import Catalog from "./catalog";
import Details from "./details";

export interface WorkflowProps extends BaseWorkflowProps {}
export interface DetailWorkflowProps {
  children?: React.ReactElement<DetailsCard>[] | React.ReactElement<DetailsCard>;
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

export { DynamicCard, MetaCard } from "./details/card";
export { LastEvent } from "./details/helpers";
export { useProjectDetailsContext } from "./details/context";
export { Details as ProjectDetails };
export type { DetailsCard, ProjectInfoChip };

export default register;
