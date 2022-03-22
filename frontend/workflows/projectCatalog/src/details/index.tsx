import React from "react";
import { useParams } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
import type { ClutchError } from "@clutch-sh/core";
import { Card, client, Grid, styled } from "@clutch-sh/core";
import { pick } from "lodash";

import type { WorkflowProps } from "..";

import type { AlertInfo } from "./alerts/types";
import type { DeployInfo, Statuses } from "./deploys/types";
import ProjectAlerts from "./alerts";
import ProjectDeploys from "./deploys";
import ProjectHeader from "./header";
import ProjectInfo from "./info";
// import ProjectContainers from "./containers";

const StyledCard = styled(Card)({
  width: "100%",
  height: "100%",
});

const StyledContainer = styled(Grid)({
  padding: "20px",
});

const StyledRow = styled(Grid)({
  minHeight: "200px",
  height: "100%",
  marginBottom: "16px",
});

interface Project {
  name?: string;
  description?: string;
  owners?: string[];
  slack?: string;
  repository: string;
  // eslint-disable-next-line camelcase
  pagerduty_escalation_policy?: string;
  prs?: number;
  languages?: string[];
  slo?: number;
  tier?: string;
  oncall?: IClutch.core.project.v1.IOnCall;
  disabled?: boolean;
}

const fetchProject = (project: string) => {
  return client
    .post("/v1/project/getProjects", { projects: [project] })
    .then(resp => {
      const { results = {} } = resp.data as IClutch.project.v1.GetProjectsResponse;
      return results[project] ? results[project].project : {};
    })
    .catch((err: ClutchError) => {
      throw err;
    });
};

const Details: React.FC<WorkflowProps> = () => {
  const { projectId } = useParams();
  const [projectInfo, setProjectInfo] = React.useState<Project>();

  // TODO:
  // customizable quick links via workflow?
  // how to modify deploys and alerts to be single project, adjsut header?
  // add tabs to facet table?
  // handle responsive sizing

  React.useEffect(() => {
    if (projectId) {
      const fetch = async () => {
        const project = await fetchProject(projectId);

        setProjectInfo({
          ...pick(project, ["name", "tier", "owners", "languages", "oncall"]),
          ...pick(project.data, [
            "description",
            "disabled",
            "pagerduty_escalation_policy",
            "repository",
            "slack",
            "team_email",
          ]),
        });
      };

      // eslint-disable-next-line no-console
      fetch().catch(console.error);
    }
  }, []);

  const fetchDeploys = (): Promise<DeployInfo> => {
    return new Promise((resolve, reject) => {
      client
        .post("/v1/lyftdeploys/getProjectEvents", {
          projects: [{ name: projectInfo.name, repo: projectInfo.repository }],
        })
        .then(res => {
          const deployInfo = {
            jobs: [
              {
                name: "test project",
                commitMetadata: {
                  repositoryName: "testing",
                  commits: [
                    {
                      ref: "1234",
                      message: "A commit message",
                      author: {
                        username: "jslaughter",
                        email: "jslaughter@lyft.com",
                      },
                    },
                    {
                      ref: "2345",
                      message: "Another commit message",
                      author: {
                        username: "jslaughter",
                        email: "jslaughter@lyft.com",
                      },
                    },
                  ],
                  baseRef: "1234",
                },
                status: "SUCCESS",
                timestamp: 1647388500000,
                environment: "PRODUCTION",
              },
              {
                name: "test 2",
                commitMetadata: {
                  repositoryName: "testing",
                  commits: [
                    {
                      ref: "1234",
                      message: "A commit message",
                      author: {
                        username: "jslaughter",
                        email: "jslaughter@lyft.com",
                      },
                    },
                    {
                      ref: "2345",
                      message: "Another commit message",
                      author: {
                        username: "jslaughter",
                        email: "jslaughter@lyft.com",
                      },
                    },
                  ],
                  baseRef: "1234",
                },
                status: "RUNNING" as Statuses,
                timestamp: 1647464082248,
                environment: "STAGING",
              },
              {
                name: "stuff",
                commitMetadata: {
                  repositoryName: "testing",
                  commits: [
                    {
                      ref: "1234",
                      message: "A commit message",
                      author: {
                        username: "jslaughter",
                        email: "jslaughter@lyft.com",
                      },
                    },
                    {
                      ref: "2345",
                      message: "Another commit message",
                      author: {
                        username: "jslaughter",
                        email: "jslaughter@lyft.com",
                      },
                    },
                  ],
                  baseRef: "1234",
                },
                status: "FAILURE" as Statuses,
                timestamp: 1647464082248,
                environment: "STAGING",
              },
            ],
            inProgress: 1,
            failures: 1,
            lastDeploy: 1647388500000,
          };

          resolve(deployInfo as DeployInfo);
        })
        .catch(reject);
    });
  };

  const fetchAlerts = (): Promise<AlertInfo> => {
    return new Promise((resolve, reject) => {
      client
        .post("/v1/lyftdeploys/getProjectEvents", {
          projects: [{ name: projectInfo.name, repo: projectInfo.repository }],
        })
        .then(res => {
          const alertInfo = {
            lastAlert: 1647388500000,
            acknowledged: 4,
            open: 2,
            projectAlerts: {
              clutch: {
                incidents: [
                  {
                    id: "1234",
                    status: "ack",
                    urgency: "HIGH",
                    url: "https://lyft.pagerduty.com/incidents/Q1U59N4B7MF6TA",
                    created: "2022-03-18T21:50:04Z",
                    description:
                      "[M3] request-ride / Failures / 4xx Details / Zeroes in ride requests to create ride segment endpoint",
                    assignments: [
                      {
                        assignee: "Josh Slaughter",
                        at: "2022-03-18T21:50:04Z",
                      },
                    ],
                  },
                  {
                    id: "1234",
                    status: "ack",
                    urgency: "LOW",
                    created: "2022-03-18T21:50:04Z",
                    description:
                      "[M3] request-ride / Failures / 4xx Details / Zeroes in ride requests to create ride segment endpoint",
                    assignments: [
                      {
                        assignee: "Josh Slaughter",
                        at: "2022-03-18T21:50:04Z",
                      },
                    ],
                  },
                  {
                    id: "1234",
                    status: "open",
                    urgency: "LOW",
                    created: "2022-03-18T21:50:04Z",
                    description:
                      "[M3] request-ride / Failures / 4xx Details / Zeroes in ride requests to create ride segment endpoint",
                    assignments: [],
                  },
                ],
                acknowledged: 2,
                open: 1,
              },
              lns: {
                incidents: [
                  {
                    id: "1234",
                    status: "ack",
                    urgency: "HIGH",
                    url: "https://lyft.pagerduty.com/incidents/Q1U59N4B7MF6TA",
                    created: "2022-03-18T21:50:04Z",
                    description:
                      "[M3] request-ride / Failures / 4xx Details / Zeroes in ride requests to create ride segment endpoint",
                    assignments: [
                      {
                        assignee: "Josh Slaughter",
                        at: "2022-03-18T21:50:04Z",
                      },
                    ],
                  },
                ],
              },
            },
          };

          resolve(alertInfo as AlertInfo);
        })
        .catch(reject);
    });
  };

  // const fetchContainers = (): Promise<> => {};

  return (
    <>
      <StyledContainer container>
        <Grid item xs={12} style={{ marginBottom: "22px" }}>
          {projectInfo && (
            <ProjectHeader name={projectInfo.name} description={projectInfo.description} />
          )}
        </Grid>
        <StyledRow container spacing={3}>
          <Grid item xs={12} sm={6} md={3} lg={3}>
            {projectInfo && (
              <ProjectInfo
                name={projectInfo.name}
                owners={projectInfo.owners}
                slack={projectInfo.slack}
                repository={projectInfo.repository}
                languages={projectInfo.languages}
                disabled={projectInfo.disabled}
                tier={projectInfo.tier}
                // slo={projectInfo.slo}
              />
            )}
          </Grid>
          <Grid item xs={12} sm={12} md={9} lg={9}>
            <StyledCard>Table</StyledCard>
          </Grid>
        </StyledRow>
        <StyledRow container spacing={3}>
          <Grid item xs={12} sm={12} md={6} lg={6}>
            {projectInfo && <ProjectDeploys fetchDeploys={fetchDeploys} />}
          </Grid>
          <Grid item xs={12} sm={12} md={6} lg={6}>
            {projectInfo && <ProjectAlerts fetchAlerts={fetchAlerts} />}
          </Grid>
        </StyledRow>
        {/* <StyledRow container spacing={3}> */}
        {/* <Grid item xs={12} sm={8} md={8} lg={8}>
            <StyledCard>Pods</StyledCard>
          </Grid> */}
        {/* <Grid item xs={12} sm={4} md={4} lg={4}> */}
        {/* <StyledCard>QL</StyledCard> */}
        {/* {projectInfo && <QuickLinks fetchLinks={fetchLinks} /> } */}
        {/* </Grid> */}
        {/* </StyledRow> */}
      </StyledContainer>
    </>
  );
};

export default Details;
