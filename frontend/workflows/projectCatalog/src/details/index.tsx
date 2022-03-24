import React from "react";
import { useParams } from "react-router-dom";
import type { ClutchError } from "@clutch-sh/core";
import { Grid, styled } from "@clutch-sh/core";

import type { DetailWorkflowProps } from "..";

import type { ProjectAlerts } from "./alerts/types";
import type { ProjectDeploys } from "./deploys/types";
import type { ProjectInfo } from "./info/types";
import ProjectAlertsCard from "./alerts";
import ProjectResourcesCard from "./containers";
import ProjectDeploysCard from "./deploys";
import ProjectHeader from "./header";
import ProjectInfoCard from "./info";
import QuickLinksCard from "./quick-links";
import { castArray } from "lodash";
import DynamicCard from "./cards/dynamic";

const StyledContainer = styled(Grid)({
  padding: "20px",
});

const Details: React.FC<DetailWorkflowProps> = ({ children }) => {
  const { projectId } = useParams();

  if (children) {
    React.Children.forEach(children, child => {
      if (React.isValidElement(child)) {
        console.log("My child", child);
        console.log((child as React.ReactElement<any>).type === DynamicCard);
      }
    });
  }
  // const [projectInfo, setProjectInfo] = React.useState<ProjectInfo>();
  // const [projectLoading, setProjectLoading] = React.useState<boolean>(false);
  // const [projectError, setProjectError] = React.useState<ClutchError | undefined>(undefined);
  // const [alertsData, setAlertsData] = React.useState<ProjectAlerts>(null);
  // const [alertsError, setAlertsError] = React.useState<ClutchError | undefined>(undefined);
  // const [alertsLoading, setAlertsLoading] = React.useState<boolean>(false);
  // const [deploysData, setDeploysData] = React.useState<ProjectDeploys>(null);
  // const [deploysError, setDeploysError] = React.useState<ClutchError | undefined>(undefined);
  // const [deploysLoading, setDeploysLoading] = React.useState<boolean>(false);

  // const refreshInterval = 30000;

  // const loadProjectInfo = () => {
  //   setProjectLoading(true);
  //   projectResolver(projectId)
  //     .then(data => setProjectInfo(data))
  //     .catch((e: ClutchError) => setProjectError(e))
  //     .finally(() => {
  //       setProjectLoading(false);
  //     });
  // };

  // const loadDeploysInfo = () => {
  //   setDeploysLoading(true);

  //   deployResolver(
  //     projectId,
  //     projectInfo?.repository?.repo,
  //     `lyft/${projectInfo?.repository?.name}`
  //   )
  //     .then(data => setDeploysData(data))
  //     .catch((e: ClutchError) => setDeploysError(e))
  //     .finally(() => setDeploysLoading(false));
  // };

  // const loadAlertsInfo = () => {
  //   setAlertsLoading(true);

  //   alertsResolver(projectInfo.serviceIds)
  //     .then(data => setAlertsData(data))
  //     .catch((e: ClutchError) => {
  //       console.log(Object.keys(e));
  //       setAlertsError(e);
  //     })
  //     .finally(() => setAlertsLoading(false));
  // };

  // React.useEffect(() => {
  //   let infoInterval;
  //   if (projectId && projectResolver && typeof projectResolver === "function") {
  //     loadProjectInfo();
  //     infoInterval = setInterval(loadProjectInfo, refreshInterval);
  //   }

  //   return () => {
  //     if (infoInterval) {
  //       clearInterval(infoInterval);
  //     }
  //   };
  // }, []);

  // React.useEffect(() => {
  //   let alertsInterval;
  //   let deploysInterval;
  //   if (projectInfo) {
  //     if (alertsResolver && typeof alertsResolver === "function") {
  //       loadAlertsInfo();
  //       alertsInterval = setInterval(loadAlertsInfo, refreshInterval);
  //     }

  //     if (deployResolver && typeof deployResolver === "function") {
  //       loadDeploysInfo();
  //       deploysInterval = setInterval(loadDeploysInfo, refreshInterval);
  //     }
  //   }

  //   return () => {
  //     if (alertsInterval) {
  //       clearInterval(alertsInterval);
  //     }

  //     if (deploysInterval) {
  //       clearInterval(deploysInterval);
  //     }
  //   };
  // }, [projectInfo]);

  return (
    <>
      <StyledContainer container>
        <Grid container direction="row">
          <Grid container item xs={11}>
            <Grid container direction="row" xs={12}>
              <Grid item style={{ marginBottom: "22px" }}>
                {/* {projectInfo && (
                  <ProjectHeader name={projectInfo.name} description={projectInfo.description} />
                )} */}
              </Grid>
              <Grid container spacing={3}>
                <Grid container item direction="row" xs={4} spacing={2}>
                  {children}
                  {/* {projectInfo && (
                    <Grid item xs={12}>
                      <ProjectInfoCard
                        data={projectInfo}
                        loading={projectLoading}
                        error={projectError}
                      />
                    </Grid>
                  )}
                  {projectInfo && (
                    <Grid item xs={12}>
                      <ProjectAlertsCard
                        data={alertsData}
                        loading={alertsLoading}
                        error={alertsError}
                      />
                    </Grid>
                  )}
                  {projectInfo && (
                    <Grid item xs={12}>
                      <ProjectDeploysCard
                        data={deploysData}
                        loading={deploysLoading}
                        error={deploysError}
                      />
                    </Grid>
                  )} */}
                </Grid>
                <Grid container item direction="row" xs={8}>
                  <ProjectResourcesCard />
                </Grid>
              </Grid>
            </Grid>
          </Grid>
          <Grid container item xs={1}>
            <QuickLinksCard />
          </Grid>
        </Grid>
      </StyledContainer>
    </>
  );
};

export default Details;
