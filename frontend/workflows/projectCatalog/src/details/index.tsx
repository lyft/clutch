import React from "react";
import { useParams } from "react-router-dom";
import { Grid, styled, Tooltip } from "@clutch-sh/core";
import { faLock } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import GroupIcon from "@material-ui/icons/Group";
import { capitalize } from "lodash";

import type { DetailWorkflowProps } from "..";

import MetaCard from "./cards/meta";
import type { ProjectInfo } from "./info/types";
import fetchProject from "./resolvers/info";
// import QuickLinksCard from "./quick-links";
import { ProjectDetailsContext } from "./context";
import ProjectHeader from "./header";
import ProjectInfoCard from "./info";

const StyledContainer = styled(Grid)({
  padding: "20px",
});

const DisabledItem = ({ name }: { name: string }) => (
  <Grid item>
    <Tooltip title={`${name} is disabled`}>
      <FontAwesomeIcon icon={faLock} size="lg" />
    </Tooltip>
  </Grid>
);

const Details: React.FC<DetailWorkflowProps> = ({ children }) => {
  const { projectId } = useParams();
  const [projectInfo, setProjectInfo] = React.useState<ProjectInfo | null>(null);
  const [metaCards, setMetaCards] = React.useState<any[]>([]);
  const [dynamicCards, setDynamicCards] = React.useState<any[]>([]);

  React.useEffect(() => {
    if (children) {
      React.Children.forEach(children, child => {
        if (React.isValidElement(child)) {
          const { type } = child?.props;

          switch (type) {
            case "Metadata":
              setMetaCards([...metaCards, child]);
              break;
            case "Dynamic":
              setDynamicCards([...dynamicCards, child]);
              break;
            default: // Do nothing, invalid card
          }
        }
      });
    }
  }, [children]);

  return (
    <>
      <ProjectDetailsContext.Provider value={{ projectInfo }}>
        <StyledContainer container>
          <Grid container direction="row">
            <Grid container item xs={11}>
              <Grid container direction="row" xs={12}>
                <Grid item style={{ marginBottom: "22px" }}>
                  {projectInfo && (
                    <ProjectHeader name={projectInfo.name} description={projectInfo.description} />
                  )}
                </Grid>
                <Grid container spacing={3}>
                  <Grid container item direction="row" xs={4} spacing={2}>
                    <Grid item xs={12}>
                      <MetaCard
                        title={capitalize(projectInfo?.name)}
                        titleIcon={<GroupIcon />}
                        fetchDataFn={() => fetchProject(projectId)}
                        onSuccess={setProjectInfo}
                        autoReload
                        endAdornment={
                          projectInfo?.disabled ? (
                            <DisabledItem name={capitalize(projectInfo?.name)} />
                          ) : null
                        }
                      >
                        {projectInfo && <ProjectInfoCard data={projectInfo} />}
                      </MetaCard>
                    </Grid>
                    {metaCards.length > 0 &&
                      metaCards.map(card => (
                        <Grid item xs={12}>
                          {card}
                        </Grid>
                      ))}
                  </Grid>
                  <Grid container item direction="row" xs={8}>
                    {dynamicCards.length > 0 &&
                      dynamicCards.map(card => (
                        <Grid item xs={12}>
                          {card}
                        </Grid>
                      ))}
                  </Grid>
                </Grid>
              </Grid>
            </Grid>
            {/* <Grid container item xs={1}>
            <QuickLinksCard />
          </Grid> */}
          </Grid>
        </StyledContainer>
      </ProjectDetailsContext.Provider>
    </>
  );
};

export default Details;
