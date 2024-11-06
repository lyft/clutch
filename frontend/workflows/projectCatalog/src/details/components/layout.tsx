import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Grid, QuickLinkGroup, useNavigate, useParams, useTheme } from "@clutch-sh/core";

import type { ProjectDetailsWorkflowProps } from "../../types";
import { ProjectDetailsContext } from "../context";
import fetchProjectInfo from "../resolver";

import type { BreadCrumbsProps } from "./breadcrumbs";
import BreadCrumbs from "./breadcrumbs";
import ProjectHeader, { ProjectHeaderProps } from "./header";
import QuickLinksAndSettings from "./link-settings";

export interface CatalogLayoutProps
  extends BreadCrumbsProps,
    ProjectHeaderProps,
    Pick<ProjectDetailsWorkflowProps, "configLinks" | "allowDisabled"> {
  children?: React.ReactNode;
  quickLinkSettings?: boolean;
}

const CatalogLayout = ({
  routes = [],
  title,
  description,
  configLinks = [],
  children,
  allowDisabled,
  quickLinkSettings = true,
}: CatalogLayoutProps) => {
  const { projectId = "" } = useParams();
  const navigate = useNavigate();
  const [projectInfo, setProjectInfo] = React.useState<IClutch.core.project.v1.IProject | null>(
    null
  );
  const projInfo = React.useMemo(() => ({ projectId, projectInfo }), [projectId, projectInfo]);
  const theme = useTheme();

  const redirectNotFound = () => navigate(`/${projectId}/notFound`, { replace: true });

  const fetchData = () =>
    fetchProjectInfo(projectId, allowDisabled)
      .then(data => {
        if (!data) {
          redirectNotFound();
          return;
        }
        setProjectInfo(data as IClutch.core.project.v1.IProject);
      })
      .catch(err => {
        // eslint-disable-next-line no-console
        console.error(err);
      });

  React.useEffect(() => {
    if (projectId) {
      fetchData();

      const interval = setInterval(fetchData, 30000);

      return () => (interval ? clearInterval(interval) : undefined);
    }

    return () => {};
  }, [projectId]);

  return (
    <ProjectDetailsContext.Provider value={projInfo}>
      <Grid container padding={theme.clutch.layout.gutter}>
        {!theme.clutch.useWorkflowLayout && (
          <Grid container item direction="column">
            <Grid item>
              <BreadCrumbs routes={[{ title: projectId, path: `${projectId}` }, ...routes]} />
            </Grid>
          </Grid>
        )}
        <Grid container item spacing={1}>
          <Grid
            container
            item
            justifyContent="space-between"
            alignItems="center"
            marginBottom="16px"
          >
            <Grid item>
              <ProjectHeader
                title={`${projectInfo?.name ?? projectId}${title ? ` ${title}` : ""}`}
                description={description ?? (projectInfo?.data?.description as string)}
              />
            </Grid>
            <Grid item>
              {projectInfo && (
                <QuickLinksAndSettings
                  linkGroups={(projectInfo.linkGroups as QuickLinkGroup[]) || []}
                  configLinks={configLinks ?? []}
                  showSettings={quickLinkSettings}
                />
              )}
            </Grid>
          </Grid>
        </Grid>
        <Grid container item spacing={2}>
          {children && children}
        </Grid>
      </Grid>
    </ProjectDetailsContext.Provider>
  );
};

export default CatalogLayout;
