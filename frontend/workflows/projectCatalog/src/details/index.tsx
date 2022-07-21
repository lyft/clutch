import React from "react";
import { useNavigate, useParams } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
import { FeatureOn, Grid, IconButton, SimpleFeatureFlag, styled, Tooltip } from "@clutch-sh/core";
import { faLock } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import GroupIcon from "@mui/icons-material/Group";
import SettingsIcon from "@mui/icons-material/Settings";
import { capitalize, isEmpty } from "lodash";

import type { CatalogDetailsChild, ProjectDetailsWorkflowProps } from "..";

import { CardType, DynamicCard, MetaCard } from "./card";
import { ProjectDetailsContext } from "./context";
import ProjectHeader from "./header";
import { fetchProjectInfo } from "./helpers";
import ProjectInfoCard from "./info";
import QuickLinksCard from "./quick-links";

const StyledContainer = styled(Grid)({
  padding: "16px",
});

const StyledHeadingContainer = styled(Grid)({
  marginBottom: "24px",
});

const DisabledItem = ({ name }: { name: string }) => (
  <Grid item>
    <Tooltip title={`${capitalize(name)} is disabled`}>
      <FontAwesomeIcon icon={faLock} size="lg" />
    </Tooltip>
  </Grid>
);

const QuickLinksAndSettingsBtn = ({ linkGroups }) => {
  const navigate = useNavigate();

  return (
    <Grid
      container
      direction="row"
      alignItems="center"
      justifyContent="flex-end"
      spacing={1}
      style={{
        padding: "8px 32px 0px 0px",
      }}
    >
      <Grid item>
        <QuickLinksCard linkGroups={linkGroups} />
      </Grid>
      <SimpleFeatureFlag feature="projectCatalogSettings">
        <FeatureOn>
          <Grid item>
            <IconButton onClick={() => navigate("config")} size="medium">
              <SettingsIcon />
            </IconButton>
          </Grid>
        </FeatureOn>
      </SimpleFeatureFlag>
    </Grid>
  );
};

const Details: React.FC<ProjectDetailsWorkflowProps> = ({ children, chips }) => {
  const { projectId } = useParams();
  const [projectInfo, setProjectInfo] = React.useState<IClutch.core.project.v1.IProject | null>(
    null
  );
  const [metaCards, setMetaCards] = React.useState<CatalogDetailsChild[]>([]);
  const [dynamicCards, setDynamicCards] = React.useState<CatalogDetailsChild[]>([]);
  const projInfo = React.useMemo(() => ({ projectInfo }), [projectInfo]);

  React.useEffect(() => {
    if (children) {
      const tempMetaCards: CatalogDetailsChild[] = [];
      const tempDynamicCards: CatalogDetailsChild[] = [];

      React.Children.forEach(children, child => {
        if (React.isValidElement(child)) {
          const { type } = child?.props || {};

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
    <ProjectDetailsContext.Provider value={projInfo}>
      <StyledContainer container direction="row" wrap="nowrap">
        {/* Column for project details and header */}
        <Grid container item direction="column" xs={12} sm={12} md={12} lg={12} xl={12}>
          <Grid container item>
            <StyledHeadingContainer item xs={6} sm={6} md={7} lg={8} xl={9}>
              {/* Static Header */}
              <ProjectHeader
                title={projectId}
                routes={[{ title: "Details" }]}
                description={projectInfo?.data?.description as string}
              />
            </StyledHeadingContainer>
            {projectInfo && !isEmpty(projectInfo?.linkGroups) && (
              <Grid container item xs={12} sm={12} md={5} lg={4} xl={3} spacing={2}>
                <QuickLinksAndSettingsBtn linkGroups={projectInfo.linkGroups} />
              </Grid>
            )}
          </Grid>
          <Grid container item direction="row" spacing={1}>
            <Grid container item xs={12} sm={12} md={5} lg={4} xl={3} spacing={2}>
              <Grid item xs={12}>
                {/* Static Info Card */}
                <MetaCard
                  title={getOwner(projectInfo?.owners ?? []) || projectId}
                  titleIcon={<GroupIcon />}
                  fetchDataFn={() => fetchProjectInfo(projectId)}
                  onSuccess={(data: unknown) =>
                    setProjectInfo(data as IClutch.core.project.v1.IProject)
                  }
                  autoReload
                  loadingIndicator={false}
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
            <Grid container item xs={12} sm={12} md={7} lg={8} xl={9} spacing={1}>
              {/* Custom Dynamic Cards */}
              {dynamicCards.length > 0 &&
                dynamicCards.map(card => (
                  <Grid item xs={12} key={`dynamic-${card.key}`}>
                    {card}
                  </Grid>
                ))}
            </Grid>
          </Grid>
        </Grid>
      </StyledContainer>
    </ProjectDetailsContext.Provider>
  );
};

export default Details;
