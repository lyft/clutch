import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";

import Catalog from "./catalog";
import Details from "./details";

export { default as DynamicCard } from "./details/cards/dynamic";
export { default as MetaCard } from "./details/cards/meta";

export type DetailsCardTypes = "Dynamic" | "Metadata";

export interface DetailsCard {
  type: DetailsCardTypes;
  title?: string;
}

export interface WorkflowProps extends BaseWorkflowProps {}
export interface DetailWorkflowProps {
  children?: React.ReactElement<DetailsCard>[] | React.ReactElement<DetailsCard>;
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

export { LastEvent, LinkText, StyledLink, StyledRow } from "./details/cards/base";
export { useProjectDetailsContext } from "./details/context";
export { Details as ProjectDetails };

export default register;
