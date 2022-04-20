import type { clutch as IClutch } from "@clutch-sh/api";
import type { ClutchError } from "@clutch-sh/core";
import { client, userId } from "@clutch-sh/core";
import _ from "lodash";

const LOCAL_STORAGE_STATE_KEY = "catalogState";

/** Attempts to load projects stored in local storage. */
const loadProjects = (): { [key: string]: string } => {
  const storedProjects = window.localStorage.getItem(LOCAL_STORAGE_STATE_KEY);
  if (storedProjects === null || storedProjects === undefined) {
    return {};
  }
  try {
    return JSON.parse(storedProjects) || {};
  } catch {
    // If stored projects are not in valid format purge them and return default projects.
    window.localStorage.removeItem(LOCAL_STORAGE_STATE_KEY);
    return {};
  }
};

const hasState = (): boolean => {
  const storedProjects = window.localStorage.getItem(LOCAL_STORAGE_STATE_KEY);
  return !(storedProjects === null || storedProjects === undefined);
};

const writeProjects = (projects: { [key: string]: string }) => {
  window.localStorage.setItem(LOCAL_STORAGE_STATE_KEY, JSON.stringify(projects));
};

const projectRequest = (includeOwned: boolean) => {
  const projects = loadProjects();
  const requestData = {
    excludeDependencies: true,
    projects: Object.keys(projects),
  } as IClutch.project.v1.GetProjectsRequest;
  if (includeOwned) {
    requestData.users = [userId()];
  }
  return requestData;
};

const fetchProjects = (
  onSuccess: (pjts: IClutch.core.project.v1.IProject[]) => void,
  onError: (e: ClutchError) => void,
  retry: number,
  includeOwned: boolean = true
) => {
  client
    .post("/v1/project/getProjects", projectRequest(includeOwned))
    .then(resp => {
      const { results } = resp.data as IClutch.project.v1.GetProjectsResponse;
      const selectedProjects = Object.values(results)
        .filter(r => r?.from?.selected)
        .map(r => r.project);

      let projectStorage = {};
      if (!Object.keys(loadProjects()).length) {
        selectedProjects.forEach(p => {
          if (p && p.name) {
            projectStorage[p.name] = Date.now();
          }
        });
        writeProjects(projectStorage);
      } else {
        projectStorage = loadProjects();
      }
      const tsProjects = selectedProjects.map(p => ({
        ...p,
        ts: p && p.name ? projectStorage[p.name] : null,
      }));
      onSuccess(_.reverse(_.sortBy(tsProjects, ["ts"])));
    })
    .catch((err: ClutchError) => {
      const projects = loadProjects();

      // will perform a regex to pull out a project name from an error message
      // sample error message -> "unable to find project: test project"
      const missingProjectMatch = err.message.match(/.*:\W(.*)/);

      if (missingProjectMatch && missingProjectMatch[1]) {
        const missingProject = missingProjectMatch[1];

        if (err.status.code === 404 && missingProject && projects?.[missingProject]) {
          delete projects?.[missingProject];
          writeProjects(projects);
          if (retry > 0) {
            fetchProjects(onSuccess, onError, retry - 1);
          }
        }
      }
      onError(err);
    });
};

const getProjects = (
  onSuccess: (pjts: IClutch.core.project.v1.IProject[]) => void,
  onError: (e: ClutchError) => void,
  includeOwned: boolean = false
) => {
  // retry count is set to the number of projects in storage since potentially
  // all projects are invalid.
  const retryCount = Object.keys(loadProjects()).length;
  fetchProjects(onSuccess, onError, retryCount, includeOwned);
};

/** Add a project to local storage. */
const addProject = (
  project: string,
  onSuccess: (pjts: IClutch.core.project.v1.IProject[]) => void,
  onError: (e: ClutchError) => void
) => {
  const storedProjects = loadProjects();
  const updatedProjects = {
    ...storedProjects,
    [project]: Date.now().toString(),
  };
  writeProjects(updatedProjects);
  getProjects(onSuccess, e => {
    writeProjects(storedProjects);
    onError(e);
  });
};

const removeProject = (
  project: string,
  onSuccess: (pjts: IClutch.core.project.v1.IProject[]) => void,
  onError: (e: ClutchError) => void
) => {
  const storedProjects = loadProjects();
  if (storedProjects?.[project] !== undefined) {
    delete storedProjects[project];
    writeProjects(storedProjects);
  }
  getProjects(onSuccess, onError, false);
};

const clearProjects = () => {
  writeProjects({});
};

export { addProject, clearProjects, getProjects, hasState, removeProject };
