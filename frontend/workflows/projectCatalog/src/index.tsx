import React from "react";
import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";

import type { AlertMassageOptions, ProjectAlerts } from "./details/alerts/types";
import fetchDeploys from "./details/deployResolver";
import type { ProjectDeploys } from "./details/deploys/types";
import type { ProjectInfo } from "./details/info/types";
import fetchAlerts from "./details/resolvers/alerts";
import fetchProject from "./details/resolvers/info";
import Catalog from "./catalog";
import Details from "./details";

export interface WorkflowProps extends BaseWorkflowProps {}
export interface DetailWorkflowProps extends WorkflowProps {
  projectResolver?: (project: string) => Promise<ProjectInfo>;
  deployResolver?: (
    project: string,
    repository: string,
    repositoryName: string
  ) => Promise<ProjectDeploys>;
  alertsResolver?: (serviceIds: string[], options?: AlertMassageOptions) => Promise<ProjectAlerts>;
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
        component: props => (
          <Details
          cards=[
            <MetaCard1 />,
            <DynamicCard2 />
          ],
            projectResolver={fetchProject}
            alertsResolver={fetchAlerts}
            deployResolver={fetchDeploys}
            {...props}
          ><MetaCard1 /><DynamicCard2 fetchFunction /></Details>
        ),
        // featureFlag: "projectCatalog",
      },
    },
  };
};

export default register;
