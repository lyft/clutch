import type { clutch as IClutch } from "@clutch-sh/api";
import type { ClutchError } from "@clutch-sh/core";
import _ from "lodash";

import {
  deriveSwitchStatus,
  exclusiveProjectDependencies,
  PROJECT_TYPE_URL,
  removeMissingProjects,
  updateGroupstate,
} from "./helpers";
import type { Action, ProjectState, State } from "./types";
import { Group } from "./types";

const selectorReducer = (state: State, action: Action): State => {
  switch (action.type) {
    // User actions.

    case "ADD_PROJECTS": {
      // a given custom project may already exist in the group so don't trigger a state update for the duplicate(s)
      const uniqueCustomProjects =
        action?.payload?.projects?.filter(
          (project: string) => !(project in state[action.payload.group])
        ) || [];
      if (uniqueCustomProjects.length === 0) {
        // since there's no new projects to add, we return the original state as it wasn't modified
        return state;
      }

      // add the custom project (currently users can only add projects to Group.Projects) to its respective group
      // and set the project state
      let newState = {
        ...state,
        [action.payload.group]: {
          ...state[action.payload.group],
          ...Object.fromEntries(
            uniqueCustomProjects.map(v => [v, { checked: true, custom: true } as ProjectState])
          ),
        },
      };

      if (action.payload.group !== Group.PROJECTS) {
        return newState;
      }

      // check to see if we have project data for project (b/c we don't need to make an API call) and if so, update the state with
      // its upstreams/downstreams
      uniqueCustomProjects.forEach(v => {
        if (!(v in newState.projectData)) {
          return;
        }

        const dependencies = newState.projectData[v]?.dependencies || {
          upstreams: {},
          downstreams: {},
        };
        const { upstreams, downstreams } = dependencies;
        upstreams?.[PROJECT_TYPE_URL]?.ids?.forEach(upstreamDep => {
          newState = updateGroupstate(newState, Group.UPSTREAM, upstreamDep, {
            checked: false,
          });
        });
        downstreams?.[PROJECT_TYPE_URL]?.ids?.forEach(downstreamDep => {
          newState = updateGroupstate(newState, Group.DOWNSTREAM, downstreamDep, {
            checked: false,
          });
        });
      });

      return newState;
    }
    case "REMOVE_PROJECTS": {
      // remove the projects from their respective group
      const newState = {
        ...state,
        [action.payload.group]: _.omit(
          state[action.payload.group],
          action?.payload?.projects || []
        ),
      };

      if (action.payload.group !== Group.PROJECTS) {
        return newState;
      }

      // for the removed projects in Group.Projects we get their upstreams/dowstreams that are exclusive
      // (i.e. not shared by other projects that weren't removed)
      const exclusiveUpstreams = [] as string[];
      const exclusiveDownstreams = [] as string[];
      action?.payload?.projects?.forEach(p => {
        /* get the upstreams/downstreams that are exclusive to the removed project(s)

        TODO: (sperry) "remove_projects" supports bulk removal & since this func checks if upstreams/downstreams are exclusive
        to a *single* project, if a dependency is shared by multiple *removed* projects, it wont't be removed. The func should
        be refactored to take in the list of removed projects to account for this scenario.
        */
        const dependencies = exclusiveProjectDependencies(state, Group.PROJECTS, p);
        exclusiveUpstreams.push(...dependencies.upstreams);
        exclusiveDownstreams.push(...dependencies.downstreams);
      });

      // remove the upstreams/downstreams exclusive to the removed project(s)
      newState[Group.UPSTREAM] = _.omit(state[Group.UPSTREAM], exclusiveUpstreams);
      newState[Group.DOWNSTREAM] = _.omit(state[Group.DOWNSTREAM], exclusiveDownstreams);

      return newState;
    }
    case "TOGGLE_PROJECTS": {
      return {
        ...state,
        [action.payload.group]: {
          ...state[action.payload.group],
          ...Object.fromEntries(
            action?.payload?.projects?.map(key => [
              key,
              {
                ...state[action.payload.group][key],
                checked: !state[action.payload.group][key].checked,
              },
            ]) || []
          ),
        },
      };
    }
    case "ONLY_PROJECTS": {
      const newState = { ...state };

      newState[action.payload.group] = Object.fromEntries(
        Object.keys(state[action.payload.group]).map(key => [
          key,
          {
            ...state[action.payload.group][key],
            checked: action?.payload?.projects?.includes(key) || false,
          },
        ])
      );

      return newState;
    }
    case "TOGGLE_ENTIRE_GROUP": {
      // the state might not match the filtered state data passed to the components, so in the func below, we should filter out the projects
      // that _were_ rendered in order to correctly evaluate the toggled checked status for all the projects in the respective group
      const filteredState = {
        ...state,
        [action.payload.group]: _.pick(
          state[action.payload.group],
          action?.payload?.projects || []
        ),
      };
      const newCheckedValue = !deriveSwitchStatus(filteredState, action.payload.group);
      return {
        ...state,
        [action.payload.group]: {
          ...state[action.payload.group],
          ...Object.fromEntries(
            action?.payload?.projects?.map(key => [
              key,
              {
                ...state[action.payload.group][key],
                checked: newCheckedValue,
              },
            ]) || []
          ),
        },
      };
    }
    // Background actions.

    case "HYDRATE_START": {
      return { ...state, loading: true };
    }
    case "HYDRATE_END": {
      let newState = { ...state, loading: false, error: undefined } as State;

      _.forIn(
        action?.payload?.result as { [k: string]: IClutch.project.v1.IProjectResult } | null,
        (v: IClutch.project.v1.IProjectResult, k: string) => {
          // user owned project vs custom project
          if (v?.from?.users && v.from.users.length > 0) {
            newState = updateGroupstate(newState, Group.PROJECTS, k, { checked: true });
          } else if (v?.from?.selected) {
            newState = updateGroupstate(newState, Group.PROJECTS, k, {
              checked: true,
              custom: true,
            });
          }

          // add each upstream/downstream for the selected or user project
          if ((v?.from?.users && v.from.users.length > 0) || v?.from?.selected) {
            const dependencies = v?.project?.dependencies || { upstreams: {}, downstreams: {} };
            const { upstreams, downstreams } = dependencies;
            upstreams?.[PROJECT_TYPE_URL]?.ids?.forEach(upstreamDep => {
              newState = updateGroupstate(newState, Group.UPSTREAM, upstreamDep, {
                checked: false,
              });
            });
            downstreams?.[PROJECT_TYPE_URL]?.ids?.forEach(downstreamDep => {
              newState = updateGroupstate(newState, Group.DOWNSTREAM, downstreamDep, {
                checked: false,
              });
            });
          }

          // stores the raw project data for each project in the API result
          if (v?.project) {
            newState.projectData[k] = v.project;
          }
        }
      );
      return newState;
    }
    case "HYDRATE_ERROR": {
      /**
       * TODO: do we want to handle the error state differently? For example, when we render the error on the UI,
       * it won't disappear unless there's a successful API call or if the user refreshes the page. If a user performs other
       * actions, such as use the toggle/checkbox/ etc. the error message will be still be on the page
       */

      const { result } = action?.payload || { result: {} };
      const validState = { ...state };

      // Will handle given partial failures
      if (result?.partialFailures) {
        validState.projectErrors = [];
        result.partialFailures.forEach(failure => {
          if (failure.details && failure.details.length) {
            validState.projectErrors.push({
              message: failure.message,
              details: failure.details[0],
            });
          }
        });

        return removeMissingProjects(
          validState,
          validState.projectErrors.map((e: any) => e.details?.name || "")
        );
      }
      if (result?.missing && result?.missing.length) {
        return removeMissingProjects(state, result?.missing || []);
      }

      const err = result as ClutchError;
      if (err.status.code !== 404) {
        return { ...state, loading: false, error: action?.payload?.result };
      }
      const errorMsg = err.message;
      let projects = validState[Group.PROJECTS];
      if (errorMsg) {
        // this message is guaranteed to follow this format from the server
        const ERROR_FMT = /unable to find project: (.+)/;
        const invalidProject = errorMsg.match(ERROR_FMT)?.[1];
        if (!invalidProject) {
          return { ...state, loading: false, error: action?.payload?.result };
        }
        const validProjects = _.difference(Object.keys(validState[Group.PROJECTS]), [
          unescape(invalidProject),
        ]);
        projects = _.pick(validState[Group.PROJECTS], validProjects);
      }
      return {
        ...validState,
        [Group.PROJECTS]: projects,
        loading: false,
        error: action?.payload?.result,
      };
    }
    default:
      throw new Error(`unknown resolver action`);
  }
};

export default selectorReducer;
