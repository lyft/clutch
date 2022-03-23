import React from "react";
import { useParams } from "react-router-dom";
// import type { clutch as IClutch } from "@clutch-sh/api";
// import type { ClutchError } from "@clutch-sh/core";
import { Grid, styled } from "@clutch-sh/core";
import { faGithub, faSlack } from "@fortawesome/free-brands-svg-icons";
import SecurityIcon from "@material-ui/icons/Security";

import type { WorkflowProps } from "..";

import type { ProjectAlerts } from "./alerts/types";
import type { ProjectInfo } from "./info/types";
import ProjectAlertsCard from "./alerts";
import ProjectResourcesCard from "./containers";
import ProjectDeploysCard, { ProjectDeploys } from "./deploys";
import ProjectHeader from "./header";
import ProjectInfoCard from "./info";
import QuickLinksCard from "./quick-links";

const StyledContainer = styled(Grid)({
  padding: "20px",
});

// interface Project {
//   name?: string;
//   description?: string;
//   owners?: string[];
//   slack?: string;
//   repository: string;
//   // eslint-disable-next-line camelcase
//   pagerduty_escalation_policy?: string;
//   prs?: number;
//   languages?: string[];
//   slo?: number;
//   tier?: string;
//   oncall?: IClutch.core.project.v1.IOnCall;
//   disabled?: boolean;
// }

// const fetchProject = (project: string) => {
//   return client
//     .post("/v1/project/getProjects", { projects: [project] })
//     .then(resp => {
//       const { results = {} } = resp.data as IClutch.project.v1.GetProjectsResponse;
//       return results[project] ? results[project].project : {};
//     })
//     .catch((err: ClutchError) => {
//       throw err;
//     });
// };

const Details: React.FC<WorkflowProps> = () => {
  const { projectId } = useParams();
  const [projectInfo, setProjectInfo] = React.useState<ProjectInfo>();
  const [projectLoading, setProjectLoading] = React.useState<boolean>(false);
  const [alertsData, setAlertsData] = React.useState<ProjectAlerts>(null);
  // const [alertsError, setAlertsError] = React.useState<ClutchError | undefined>(undefined);
  // const [alertsLoading, setAlertsLoading] = React.useState<boolean>(false);
  const [deploysData, setDeploysData] = React.useState<ProjectDeploys>(null);
  // const [deploysError, setDeploysError] = React.useState<ClutchError | undefined>(undefined);
  // const [deploysLoading, setDeploysLoading] = React.useState<boolean>(false);

  // TODO:
  // customizable quick links via workflow?
  // how to modify deploys and alerts to be single project, adjsut header?
  // add tabs to facet table?
  // handle responsive sizing

  React.useEffect(() => {
    if (projectId) {
      const fetch = async () => {
        // const project = await fetchProject(projectId);

        // setProjectInfo({
        //   ...pick(project, ["name", "tier", "owners", "languages", "oncall"]),
        //   ...pick(project.data, [
        //     "description",
        //     "disabled",
        //     "pagerduty_escalation_policy",
        //     "repository",
        //     "slack",
        //     "team_email",
        //   ]),
        // });
        setProjectLoading(true);
        setProjectInfo({
          owner: "Platform Tools",
          name: projectId,
          disabled: false,
          description: "Clutch stuff",
          repository: {
            name: "clutch-private",
            icon: faGithub,
            url: "https://lyft.slack.com/app_redirect?channel=platform-tools",
            requests: {
              number: 11,
              type: "Open",
              url: "https://lyft.slack.com/app_redirect?channel=platform-tools",
            },
          },
          languages: ["go", "python"],
          messenger: {
            text: "#platform-tools",
            icon: faSlack,
            url: "https://lyft.slack.com/app_redirect?channel=platform-tools",
          },
          chips: [
            {
              text: "T0",
              title: "Tier 0 Service",
            },
            {
              text: "SLO 94%",
              title: "SLO Score of 94%",
              variant: "error",
              url: "https://google.com",
            },
            {
              text: "97%",
              title: "Windshield Score of 97%",
              variant: "active",
              icon: <SecurityIcon />,
            },
          ],
        });

        setAlertsData({
          title: "Pager Duty Alerts",
          lastAlert: new Date().getTime(),
          summary: {
            open: {
              count: 4,
              url: "https://www.google.com",
            },
            triggered: {
              count: 2,
            },
            acknowledged: {
              count: 23,
            },
          },
          onCall: {
            text: "Slack to page Oncall",
            icon: faSlack,
            users: [
              {
                name: "Derek Schaller",
              },
              {
                name: "Daniel Hochman",
                url: "test",
              },
            ],
            url: "stuff",
          },
          create: {
            text: "Create an Incident",
            url: "https://www.google.com",
          },
        });
      };

      setDeploysData({
        lastDeploy: new Date().getTime(),
        seeMore: {
          text: "See More Deploys",
          url: "http://www.google.com",
        },
        deploys: [
          {
            repositoryName: "clutch-private",
            commits: [
              {
                ref: "6cd92df18e53ba416131f99f17d7dbd19ad16ef5",
                message: "Update submodule",
                author: {
                  username: "GitHub CI",
                  email: "platform-tools@lyft.com",
                },
              },
              {
                ref: "f656c3a69bda8ce382061a6bf3b3323a2ca3b4d2",
                message: "Update datum",
                author: {
                  username: "GitHub CI",
                  email: "platform-tools@lyft.com",
                },
              },
              {
                ref: "af97d74f9b42ec7aa3c6b49e1cbf1db138783fb3",
                message: "Update submodule",
                author: {
                  username: "GitHub CI",
                  email: "platform-tools@lyft.com",
                },
              },
              {
                ref: "c50aa79fbd2f318adb684a4ed37a67a18fa4e869",
                message: "todos (#1375)\n",
                author: {
                  username: "Shawna Monero",
                  email: "66325812+smonero@users.noreply.github.com",
                },
              },
              {
                ref: "e39e8c629ab7837426e571c5ce7e7250a7922e54",
                message: "Update submodule",
                author: {
                  username: "GitHub CI",
                  email: "platform-tools@lyft.com",
                },
              },
              {
                ref: "60d1d5009d859a602c4f8df582a631a024a00249",
                message: "Update submodule",
                author: {
                  username: "GitHub CI",
                  email: "platform-tools@lyft.com",
                },
              },
              {
                ref: "85ea92fc95c1bf97b0df2a4af9f951503989a3e7",
                message: "Update submodule",
                author: {
                  username: "GitHub CI",
                  email: "platform-tools@lyft.com",
                },
              },
              {
                ref: "1d19591688d97dfb5c2a5c401ac3f0b5b6a632fa",
                message: "Update submodule",
                author: {
                  username: "GitHub CI",
                  email: "platform-tools@lyft.com",
                },
              },
              {
                ref: "12fb06642846d442dc976a2d001081264b3cc3d7",
                message: "project api: support excludeDependencies option (#1377)\n",
                author: {
                  username: "Mike Cutalo",
                  email: "mcutalo88@gmail.com",
                },
              },
            ],
            environment: "PRODUCTION",
            baseRef: "1114",
          },
          {
            repositoryName: "clutch-private",
            commits: [
              {
                ref: "6cd92df18e53ba416131f99f17d7dbd19ad16ef5",
                message: "Update submodule",
                author: {
                  username: "GitHub CI",
                  email: "platform-tools@lyft.com",
                },
              },
              {
                ref: "f656c3a69bda8ce382061a6bf3b3323a2ca3b4d2",
                message: "Update datum",
                author: {
                  username: "GitHub CI",
                  email: "platform-tools@lyft.com",
                },
              },
              {
                ref: "af97d74f9b42ec7aa3c6b49e1cbf1db138783fb3",
                message: "Update submodule",
                author: {
                  username: "GitHub CI",
                  email: "platform-tools@lyft.com",
                },
              },
              {
                ref: "c50aa79fbd2f318adb684a4ed37a67a18fa4e869",
                message: "todos (#1375)\n",
                author: {
                  username: "Shawna Monero",
                  email: "66325812+smonero@users.noreply.github.com",
                },
              },
              {
                ref: "e39e8c629ab7837426e571c5ce7e7250a7922e54",
                message: "Update submodule",
                author: {
                  username: "GitHub CI",
                  email: "platform-tools@lyft.com",
                },
              },
              {
                ref: "60d1d5009d859a602c4f8df582a631a024a00249",
                message: "Update submodule",
                author: {
                  username: "GitHub CI",
                  email: "platform-tools@lyft.com",
                },
              },
              {
                ref: "85ea92fc95c1bf97b0df2a4af9f951503989a3e7",
                message: "Update submodule",
                author: {
                  username: "GitHub CI",
                  email: "platform-tools@lyft.com",
                },
              },
              {
                ref: "1d19591688d97dfb5c2a5c401ac3f0b5b6a632fa",
                message: "Update submodule",
                author: {
                  username: "GitHub CI",
                  email: "platform-tools@lyft.com",
                },
              },
              {
                ref: "12fb06642846d442dc976a2d001081264b3cc3d7",
                message: "project api: support excludeDependencies option (#1377)\n",
                author: {
                  username: "Mike Cutalo",
                  email: "mcutalo88@gmail.com",
                },
              },
            ],
            environment: "STAGING",
            baseRef: "1114",
          },
        ],
      });

      // eslint-disable-next-line no-console
      fetch()
        .catch(console.error)
        .finally(() => setProjectLoading(false));
    }
  }, []);

  // const fetchDeploys = () => {
  //   setDeploysLoading(true);
  //   client
  //     .post("/v1/lyftdeploys/getProjectEvents", {
  //       projects: [{ name: projectInfo.name, repo: projectInfo.repository }],
  //     })
  //     .then(res => {
  //       const deployInfo = {
  //         jobs: [
  //           {
  //             name: "test project",
  //             commitMetadata: {
  //               repositoryName: "testing",
  //               commits: [
  //                 {
  //                   ref: "1234",
  //                   message: "A commit message",
  //                   author: {
  //                     username: "jslaughter",
  //                     email: "jslaughter@lyft.com",
  //                   },
  //                 },
  //                 {
  //                   ref: "2345",
  //                   message: "Another commit message",
  //                   author: {
  //                     username: "jslaughter",
  //                     email: "jslaughter@lyft.com",
  //                   },
  //                 },
  //               ],
  //               baseRef: "1234",
  //             },
  //             status: "SUCCESS",
  //             timestamp: 1647388500000,
  //             environment: "PRODUCTION",
  //           },
  //           {
  //             name: "test 2",
  //             commitMetadata: {
  //               repositoryName: "testing",
  //               commits: [
  //                 {
  //                   ref: "1234",
  //                   message: "A commit message",
  //                   author: {
  //                     username: "jslaughter",
  //                     email: "jslaughter@lyft.com",
  //                   },
  //                 },
  //                 {
  //                   ref: "2345",
  //                   message: "Another commit message",
  //                   author: {
  //                     username: "jslaughter",
  //                     email: "jslaughter@lyft.com",
  //                   },
  //                 },
  //               ],
  //               baseRef: "1234",
  //             },
  //             status: "RUNNING" as Statuses,
  //             timestamp: 1647464082248,
  //             environment: "STAGING",
  //           },
  //           {
  //             name: "stuff",
  //             commitMetadata: {
  //               repositoryName: "testing",
  //               commits: [
  //                 {
  //                   ref: "1234",
  //                   message: "A commit message",
  //                   author: {
  //                     username: "jslaughter",
  //                     email: "jslaughter@lyft.com",
  //                   },
  //                 },
  //                 {
  //                   ref: "2345",
  //                   message: "Another commit message",
  //                   author: {
  //                     username: "jslaughter",
  //                     email: "jslaughter@lyft.com",
  //                   },
  //                 },
  //               ],
  //               baseRef: "1234",
  //             },
  //             status: "FAILURE" as Statuses,
  //             timestamp: 1647464082248,
  //             environment: "STAGING",
  //           },
  //         ],
  //         inProgress: 1,
  //         failures: 1,
  //         lastDeploy: 1647388500000,
  //       };

  //       setDeploysData(deployInfo as DeployInfo);
  //     })
  //     .catch(e => setDeploysError(e))
  //     .finally(() => setDeploysLoading(false));
  // };

  // const fetchAlerts = () => {
  //   setAlertsLoading(true);

  //   client
  //     .post("/v1/lyftdeploys/getProjectEvents", {
  //       projects: [{ name: projectInfo.name, repo: projectInfo.repository }],
  //     })
  //     .then(res => {
  //       const alertInfo = {
  //         lastAlert: 1647388500000,
  //         acknowledged: 4,
  //         open: 2,
  //         projectAlerts: {
  //           clutch: {
  //             incidents: [
  //               {
  //                 id: "1234",
  //                 status: "ack",
  //                 urgency: "HIGH",
  //                 url: "https://lyft.pagerduty.com/incidents/Q1U59N4B7MF6TA",
  //                 created: "2022-03-18T21:50:04Z",
  //                 description:
  //                   "[M3] request-ride / Failures / 4xx Details / Zeroes in ride requests to create ride segment endpoint",
  //                 assignments: [
  //                   {
  //                     assignee: "Josh Slaughter",
  //                     at: "2022-03-18T21:50:04Z",
  //                   },
  //                 ],
  //               },
  //               {
  //                 id: "1234",
  //                 status: "ack",
  //                 urgency: "LOW",
  //                 created: "2022-03-18T21:50:04Z",
  //                 description:
  //                   "[M3] request-ride / Failures / 4xx Details / Zeroes in ride requests to create ride segment endpoint",
  //                 assignments: [
  //                   {
  //                     assignee: "Josh Slaughter",
  //                     at: "2022-03-18T21:50:04Z",
  //                   },
  //                 ],
  //               },
  //               {
  //                 id: "1234",
  //                 status: "open",
  //                 urgency: "LOW",
  //                 created: "2022-03-18T21:50:04Z",
  //                 description:
  //                   "[M3] request-ride / Failures / 4xx Details / Zeroes in ride requests to create ride segment endpoint",
  //                 assignments: [],
  //               },
  //             ],
  //             acknowledged: 2,
  //             open: 1,
  //           },
  //           lns: {
  //             incidents: [
  //               {
  //                 id: "1234",
  //                 status: "ack",
  //                 urgency: "HIGH",
  //                 url: "https://lyft.pagerduty.com/incidents/Q1U59N4B7MF6TA",
  //                 created: "2022-03-18T21:50:04Z",
  //                 description:
  //                   "[M3] request-ride / Failures / 4xx Details / Zeroes in ride requests to create ride segment endpoint",
  //                 assignments: [
  //                   {
  //                     assignee: "Josh Slaughter",
  //                     at: "2022-03-18T21:50:04Z",
  //                   },
  //                 ],
  //               },
  //             ],
  //           },
  //         },
  //       };
  //       setAlertsData(alertInfo as AlertInfo);
  //     })
  //     .catch(e => setAlertsError(e))
  //     .finally(() => setAlertsLoading(false));
  // };

  // const fetchContainers = (): Promise<> => {};

  return (
    <>
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
                  {projectInfo && (
                    <Grid item xs={12}>
                      <ProjectInfoCard info={projectInfo} loading={projectLoading} />
                    </Grid>
                  )}
                  {alertsData && (
                    <Grid item xs={12}>
                      <ProjectAlertsCard {...alertsData} />
                    </Grid>
                  )}
                  {deploysData && (
                    <Grid item xs={12}>
                      <ProjectDeploysCard {...deploysData} />
                    </Grid>
                  )}
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
