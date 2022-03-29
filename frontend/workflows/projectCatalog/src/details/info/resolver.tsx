import type { clutch as IClutch } from "@clutch-sh/api";
import { client } from "@clutch-sh/core";
import { faGithub, faSlack } from "@fortawesome/free-brands-svg-icons";

import type { ProjectInfo, ProjectRepository } from "./types";

const getRepoData = (repo: string): ProjectRepository => {
  const repoInfo: any = {
    repo,
  };

  const [base, splitProject] = repo.split(":");
  const project = splitProject.replace(".git", "");
  const [, manager] = base.split("@");

  // nested project
  if (project.indexOf("/") > 0) {
    repoInfo.name = project.split("/").pop() || project;
  } else {
    repoInfo.name = project;
  }

  if (manager) {
    repoInfo.url = `https://${manager}/${project}`;

    switch (manager.toLowerCase()) {
      case "github.com":
        repoInfo.icon = faGithub;
        break;
      default: // no icon
    }
  }

  return repoInfo as ProjectRepository;
};

const massageProjectInfo = (project: IClutch.core.project.v1.IProject = {}): ProjectInfo => {
  const messengerLink = (project?.linkGroups || []).find(lg => lg.name === "Slack");

  return {
    owner: (project?.owners || [])[0] ?? "",
    name: project?.name as string,
    disabled: project?.data?.disabled as boolean,
    description: project?.data?.description as string,
    serviceIds: (project?.data?.pagerduty_service_ids || []) as string[],
    repository: getRepoData(project?.data?.repository as string),
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
  client
    .post("/v1/project/getProjects", { projects: [project], excludeDependencies: true })
    .then(resp => {
      const { results = {} } = resp.data as IClutch.project.v1.GetProjectsResponse;

      return massageProjectInfo(results[project] ? results[project].project ?? {} : {});
    });

export default fetchProject;
