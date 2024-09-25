import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Grid, QuickLinkGroup, styled, useNavigate, useParams } from "@clutch-sh/core";
import { ProjectDetailsContext } from "../context";
import BreadCrumbs from "./breadcrumbs";
import type { BreadCrumbsProps } from "./breadcrumbs";
import ProjectHeader, { ProjectHeaderProps } from "./header";
import QuickLinksAndSettings from "./link-settings";
import { ProjectDetailsWorkflowProps } from "..";
import fetchProjectInfo from "../resolver";

export interface CatalogLayoutProps
  extends BreadCrumbsProps,
    ProjectHeaderProps,
    Pick<ProjectDetailsWorkflowProps, "configLinks" | "allowDisabled"> {
  children?: React.ReactNode;
  quickLinkSettings?: boolean;
}

const StyledContainer = styled(Grid)({
  padding: "8px 24px",
});

const CatalogLayout = ({
  routes = [],
  title,
  description,
  configLinks = [],
  children,
  allowDisabled,
  quickLinkSettings = true,
}: CatalogLayoutProps) => {
  const { projectId } = useParams();
  const navigate = useNavigate();
  const [projectInfo, setProjectInfo] = React.useState<IClutch.core.project.v1.IProject | null>(
    null
  );
  const projInfo = React.useMemo(() => ({ projectId, projectInfo }), [projectId, projectInfo]);

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
        console.error(err);
      });

  React.useEffect(() => {
    let interval;

    fetchData();

    interval = setInterval(fetchData, 30000);

    return () => (interval ? clearInterval(interval) : undefined);
  }, []);

  return (
    <ProjectDetailsContext.Provider value={projInfo}>
      <StyledContainer container>
        <Grid container item direction="column">
          <Grid item>
            <BreadCrumbs routes={[{ title: projectId, path: `${projectId}` }, ...routes]} />
          </Grid>
        </Grid>
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
      </StyledContainer>
    </ProjectDetailsContext.Provider>
  );
};

export default CatalogLayout;
