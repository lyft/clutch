import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";

import type DynamicCard from "./details/cards/dynamic";
import type MetaCard from "./details/cards/meta";
import Catalog from "./catalog";
import Details from "./details";

export type DetailsCardTypes = "Dynamic" | "Metadata";

export interface DetailsCard {
  type: DetailsCardTypes;
  title?: string;
}

export interface WorkflowProps extends BaseWorkflowProps {}
export interface DetailWorkflowProps {
  children?:
    | React.ReactElement<DetailsCard>[]
    | React.ReactElement<DetailsCard>
    | (DynamicCard | MetaCard)[]
    | (DynamicCard | MetaCard);
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

export { Details as ProjectDetails };

export default register;
