import type { clutch as IClutch } from "@clutch-sh/api";
import { client } from "@clutch-sh/core";

import type { ProjectDeploys } from "../deploys/types";

const massageProjectDeploys = (
  project: string,
  repositoryName: string,
  deployEvents: IClutch.timeseries.v1.IPoint[]
): ProjectDeploys => {
  return {} as ProjectDeploys;
};

const fetchDeploys = (
  project: string,
  repository: string,
  repositoryName: string
): Promise<ProjectDeploys> =>
  client
    .post("/v1/DEPLOYSURL", {
      projects: [{ name: project, repo: repository }],
      pageSize: 10,
    })
    .then(res => {
      const { data } = res;

      return massageProjectDeploys(project, repositoryName, data ?? {});
    });

export default fetchDeploys;
