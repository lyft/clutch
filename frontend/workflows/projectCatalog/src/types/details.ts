import type { DetailCard } from "../details/components/card";
import type { ProjectInfoChip } from "../details/info/chipsRow";

import type { ProjectConfigLink } from "./config";
import type { ProjectCatalogProps, WorkflowProps } from "./workflow";

export interface ProjectDetailsWorkflowProps extends WorkflowProps, ProjectCatalogProps {
  children?:
    | ((CatalogDetailsChild | CatalogDetailsChild[]) &
        (React.ReactChild | React.ReactFragment | React.ReactPortal | null))
    | undefined;
  chips?: ProjectInfoChip[];
  configLinks?: ProjectConfigLink[];
}

export type CatalogDetailsChild = React.ReactElement<DetailCard>;
