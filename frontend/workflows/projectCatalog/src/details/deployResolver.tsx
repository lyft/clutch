import type { clutch as IClutch } from "@clutch-sh/api";
import type { lyft as ILyft } from "@clutch-sh/api-private";
import { client } from "@clutch-sh/core";
import { get } from "lodash";

import type { CommitInfo, ProjectDeploys } from "./deploys/types";

const massageProjectDeploys = (
  project: string,
  repositoryName: string,
  deployEvents: IClutch.timeseries.v1.IPoint[]
): ProjectDeploys => {
  const deploys: CommitInfo[] = deployEvents.map(event => ({
    repositoryName,
    baseRef: get(event, ["pb", "commits", "0", "ref"], ""),
    url: get(event, ["pb", "jobInfo", "url"]),
    commits: get(event, ["pb", "commits"]),
    environment: get(event, ["pb", "jobInfo", "environment"]),
  }));

  return {
    title: "Deploys",
    deploys,
    lastDeploy: get(deployEvents, ["0", "pb", "startTimeMillis"]),
    seeMore: {
      text: "See More Deploys",
      url: `https://deployview.lyft.net/pipeline/${project}`,
    },
  } as ProjectDeploys;
};

const fetchDeploys = (
  project: string,
  repository: string,
  repositoryName: string
): Promise<ProjectDeploys> =>
  client
    .post("/v1/lyftdeploys/getProjectEvents", {
      projects: [{ name: project, repo: repository }],
      pageSize: 10,
    } as ILyft.deploys.v1.GetProjectEventsRequest)
    .then(res => {
      const {
        data = { projectToEventsMap: {} },
      } = res as ILyft.deploys.v1.GetProjectEventsResponse;

      return massageProjectDeploys(
        project,
        repositoryName,
        data.projectToEventsMap[project]?.events ?? {}
      );
    });

export default fetchDeploys;
