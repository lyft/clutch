import type { BaseWorkflowProps } from "@clutch-sh/core";

export interface ProjectCatalogProps {
  allowDisabled?: boolean;
}

export interface WorkflowProps extends BaseWorkflowProps, ProjectCatalogProps {}
