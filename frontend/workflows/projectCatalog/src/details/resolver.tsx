import type { clutch as IClutch } from "@clutch-sh/api";
import { client } from "@clutch-sh/core";
import * as _ from "lodash";

const fetchProjectInfo = (
  project: string,
  allowDisabled: boolean = false
): Promise<IClutch.core.project.v1.IProject | null | undefined> =>
  client
    .post("/v1/project/getProjects", {
      projects: [project],
      excludeDependencies: true,
    } as IClutch.project.v1.GetProjectsRequest)
    .then(resp => {
      const { results = {}, partialFailures } = resp.data as IClutch.project.v1.GetProjectsResponse;

      const projectResults = _.get(results, [project, "project"], {});

      if (_.isEmpty(projectResults) && allowDisabled && partialFailures) {
        // Will add disabled projects to the list of projects to display if requested
        const failuresMap = partialFailures
          .map(p => {
            if ((_.get(p, "message", "") ?? "").includes("disabled")) {
              const disabledProject = _.get(p, "details.[0]", { name: "unknown " });
              _.set(
                disabledProject,
                ["data", "description"],
                `${disabledProject.name} is disabled.`
              );
              return disabledProject;
            }
            return null;
          })
          .filter(Boolean);

        if (failuresMap.length) {
          return failuresMap[0];
        }
      }

      return projectResults as IClutch.core.project.v1.IProject;
    });

export default fetchProjectInfo;
