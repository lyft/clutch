import React from "react";
import { useParams } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
import { client, Grid, styled, Tooltip } from "@clutch-sh/core";
import { faLock } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import Hidden from "@material-ui/core/Hidden";
import GroupIcon from "@material-ui/icons/Group";
import { capitalize, isEmpty } from "lodash";

import type { CatalogDetailsChild, ProjectDetailsWorkflowProps } from "..";

import { CardType, DynamicCard, MetaCard } from "./card";
import { ProjectDetailsContext } from "./context";
import ProjectHeader from "./header";
import ProjectInfoCard from "./info";
import QuickLinksCard from "./quick-links";

const StyledContainer = styled(Grid)({
  padding: "32px",
});

const StyledHeadingContainer = styled(Grid)({
  marginBottom: "24px",
});

const StyledQLContainer = styled(Grid)({
  margin: "-8px 16px 16px 0px",
});

const DisabledItem = ({ name }: { name: string }) => (
  <Grid item>
    <Tooltip title={`${capitalize(name)} is disabled`}>
      <FontAwesomeIcon icon={faLock} size="lg" />
    </Tooltip>
  </Grid>
);

const fetchProject = (project: string): Promise<IClutch.core.project.v1.IProject> =>
  client
    .post("/v1/project/getProjects", { projects: [project], excludeDependencies: true })
    .then(resp => {
      const { results = {} } = resp.data as IClutch.project.v1.GetProjectsResponse;

      return results[project] ? results[project].project ?? {} : {};
    });

const Details: React.FC<ProjectDetailsWorkflowProps> = ({ children, chips }) => {
  const { projectId } = useParams();
  const [projectInfo, setProjectInfo] = React.useState<IClutch.core.project.v1.IProject | null>(
    null
  );
  const [metaCards, setMetaCards] = React.useState<CatalogDetailsChild[]>([]);
  const [dynamicCards, setDynamicCards] = React.useState<CatalogDetailsChild[]>([]);

  React.useEffect(() => {
    if (children) {
      const tempMetaCards: CatalogDetailsChild[] = [];
      const tempDynamicCards: CatalogDetailsChild[] = [];

      React.Children.forEach(children, child => {
        if (React.isValidElement(child)) {
          const { type } = child?.props;

          switch (type) {
            case CardType.METADATA:
              tempMetaCards.push(child);
              break;
            case CardType.DYNAMIC:
              tempDynamicCards.push(child);
              break;
            default: {
              if (child.type === DynamicCard) {
                tempDynamicCards.push(child);
              } else if (child.type === MetaCard) {
                tempMetaCards.push(child);
              }
              // Do nothing, invalid card
            }
          }
        }
      });

      setMetaCards(tempMetaCards);
      setDynamicCards(tempDynamicCards);
    }
  }, [children]);

  /**
   * Takes an array of owner emails and returns the first one capitalized
   * (ex: ["clutch-team@lyft.com"] -> "Clutch Team")
   */
  const getOwner = (owners: string[]): string => {
    if (owners && owners.length) {
      const firstOwner = owners[0];

      return firstOwner
        .replace(/@.*\..*/, "")
        .split("-")
        .join(" ");
    }

    return "";
  };

  return (
    <ProjectDetailsContext.Provider value={{ projectInfo }}>
      <StyledContainer container direction="row" wrap="nowrap">
        {/* Column for project details and header */}
        <Grid item direction="column" xs={12} sm={12} md={12} lg={12} xl={12}>
          <StyledHeadingContainer item>
            {/* Static Header */}
            <ProjectHeader
              name={projectId}
              description={projectInfo?.data?.description as string}
            />
          </StyledHeadingContainer>
          <Hidden mdUp>
            <StyledQLContainer item direction="row" xs={12} sm={12}>
              {projectInfo && <QuickLinksCard linkGroups={projectInfo?.linkGroups ?? []} />}
            </StyledQLContainer>
          </Hidden>
          <Grid container direction="row" spacing={2}>
            <Grid container item xs={12} sm={12} md={5} lg={4} xl={3} spacing={2}>
              <Grid item xs={12}>
                {/* Static Info Card */}
                <MetaCard
                  title={getOwner(projectInfo?.owners ?? []) || projectId}
                  titleIcon={<GroupIcon />}
                  fetchDataFn={() => fetchProject(projectId)}
                  onSuccess={(data: unknown) =>
                    setProjectInfo(data as IClutch.core.project.v1.IProject)
                  }
                  autoReload
                  endAdornment={
                    projectInfo?.data?.disabled ? <DisabledItem name={projectId} /> : null
                  }
                >
                  {projectInfo && <ProjectInfoCard projectData={projectInfo} addtlChips={chips} />}
                </MetaCard>
              </Grid>
              {/* Custom Meta Cards */}
              {metaCards.length > 0 &&
                metaCards.map(card => (
                  <Grid item xs={12}>
                    {card}
                  </Grid>
                ))}
            </Grid>
            <Grid container item xs={12} sm={12} md={7} lg={8} xl={9} spacing={2}>
              {/* Custom Dynamic Cards */}
              {dynamicCards.length > 0 &&
                dynamicCards.map(card => (
                  <Grid item xs={12}>
                    {card}
                  </Grid>
                ))}
            </Grid>
          </Grid>
        </Grid>
        {/* Column for project quick links */}
        <Hidden smDown>
          <Grid item direction="column">
            {projectInfo && !isEmpty(projectInfo?.linkGroups) && (
              <QuickLinksCard linkGroups={projectInfo.linkGroups!} />
            )}
          </Grid>
        </Hidden>
      </StyledContainer>
    </ProjectDetailsContext.Provider>
  );
};

export default Details;
