import React from "react";
import { useParams } from "react-router-dom";
import { Grid, styled, Tooltip } from "@clutch-sh/core";
import { faLock } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import GroupIcon from "@material-ui/icons/Group";
import { capitalize } from "lodash";

import type { DetailWorkflowProps } from "..";

import fetchProject from "./info/resolver";
import type { ProjectInfo } from "./info/types";
import type { DetailsCard } from "./card";
import { MetaCard } from "./card";
import { ProjectDetailsContext } from "./context";
import ProjectHeader from "./header";
import ProjectInfoCard from "./info";
import QuickLinksCard from "./quick-links";

type CardTyping = React.ReactNode | DetailsCard;

const StyledContainer = styled(Grid)({
  padding: "20px",
});

const DisabledItem = ({ name }: { name: string }) => (
  <Grid item>
    <Tooltip title={`${capitalize(name)} is disabled`}>
      <FontAwesomeIcon icon={faLock} size="lg" />
    </Tooltip>
  </Grid>
);

const Details: React.FC<DetailWorkflowProps> = ({ children, chips }) => {
  const { projectId } = useParams();
  const [projectInfo, setProjectInfo] = React.useState<ProjectInfo | null>(null);
  const [metaCards, setMetaCards] = React.useState<CardTyping[]>([]);
  const [dynamicCards, setDynamicCards] = React.useState<CardTyping[]>([]);

  React.useEffect(() => {
    if (children) {
      const tempMetaCards: CardTyping[] = [];
      const tempDynamicCards: CardTyping[] = [];

      React.Children.forEach(children, child => {
        if (React.isValidElement(child)) {
          const { type } = child?.props;

          switch (type) {
            case "Metadata":
              tempMetaCards.push(child);
              break;
            case "Dynamic":
              tempDynamicCards.push(child);
              break;
            default: // Do nothing, invalid card
          }
        }
      });

      setMetaCards(tempMetaCards);
      setDynamicCards(tempDynamicCards);
    }
  }, [children]);

  return (
    <ProjectDetailsContext.Provider value={{ projectInfo }}>
      <StyledContainer container direction="row" wrap="nowrap">
        {/* Column for project details and header */}
        <Grid container item direction="column">
          <Grid item style={{ marginBottom: "22px" }}>
            {/* Static Header */}
            {projectInfo && (
              <ProjectHeader name={projectInfo.name} description={projectInfo.description} />
            )}
          </Grid>
          <Grid container spacing={3}>
            <Grid container item direction="row" xs={4} spacing={2}>
              <Grid item xs={12}>
                {/* Static Info Card */}
                <MetaCard
                  title={capitalize(projectInfo?.name)}
                  titleIcon={<GroupIcon />}
                  fetchDataFn={() => fetchProject(projectId)}
                  onSuccess={setProjectInfo}
                  autoReload
                  endAdornment={
                    projectInfo?.disabled ? <DisabledItem name={projectInfo?.name} /> : null
                  }
                >
                  {projectInfo && <ProjectInfoCard data={projectInfo} addtlChips={chips} />}
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
            <Grid container item direction="row" xs={8}>
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
        <Grid container item direction="column" xs={1}>
          <QuickLinksCard />
        </Grid>
      </StyledContainer>
    </ProjectDetailsContext.Provider>
  );
};

export default Details;
