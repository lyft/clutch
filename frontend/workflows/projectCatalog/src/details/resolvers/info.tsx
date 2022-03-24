import type { clutch as IClutch } from "@clutch-sh/api";
import { client } from "@clutch-sh/core";
import { faGithub, faSlack } from "@fortawesome/free-brands-svg-icons";

import type { ProjectInfo } from "../info/types";

const massageProjectInfo = (project: IClutch.core.project.v1.IProject = {}): ProjectInfo => {
  const messengerLink = (project?.linkGroups || []).find(lg => lg.name === "Slack");

  return {
    owner: (project?.owners || [])[0] ?? "",
    name: project?.name as string,
    disabled: project?.data?.disabled as boolean,
    description: project?.data?.description as string,
    serviceIds: (project?.data?.pagerduty_service_ids || []) as string[],
    repository: {
      name: (project?.data?.repository as string).split("/")[1].split(".git")[0],
      repo: project?.data?.repository as string,
      url: "",
      icon: faGithub,
    },
    messenger: {
      text: project?.data?.slack as string,
      icon: faSlack,
      url: (messengerLink?.links || [])[0]?.url as string,
    },
    languages: project.languages ?? [],
    chips: [
      {
        text: `T${project.tier}`,
        title: `Tier ${project.tier} Service`,
      },
    ],
  } as ProjectInfo;
};

const fetchProject = (project: string): Promise<ProjectInfo> =>
  client.post("/v1/project/getProjects", { projects: [project] }).then(resp => {
    const { results = {} } = resp.data as IClutch.project.v1.GetProjectsResponse;

    return massageProjectInfo(results[project] ? results[project].project ?? {} : {});
  });

export default fetchProject;
