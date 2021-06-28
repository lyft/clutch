import type { clutch as IClutch } from "@clutch-sh/api";
import _ from "lodash";

import type { Action, ProjectState, State } from "./hello-world";
import { deriveSwitchStatus, Group } from "./hello-world";

const updateGroupstate = (
  state: State,
  group: Group,
  project: string,
  projectState: ProjectState
): State => {
  const newState = { ...state };
  if (project in newState[group]) {
    // preserve the checked value if the project is already in the group
    newState[group][project].checked = state[group][project].checked;
    newState[group][project].custom = projectState.custom;
  } else {
    // we set the projectState to the default state passed in
    newState[group][project] = projectState;
  }

  return newState;
};

const selectorReducer = (state: State, action: Action): State => {
  switch (action.type) {
    // User actions.

    case "ADD_PROJECTS": {
      // a given custom project may already exist in the group so don't trigger a state update for those duplicates
      const uniqueCustomProjects = action.payload.projects.filter(
        (project: string) => !(project in state[action.payload.group])
      );
      if (uniqueCustomProjects.length === 0) {
        return state;
      }
      return {
        ...state,
        [action.payload.group]: {
          ...state[action.payload.group],
          ...Object.fromEntries(
            uniqueCustomProjects.map(v => [v, { checked: true, custom: true }])
          ),
        },
      };
    }
    case "REMOVE_PROJECTS": {
      const newRemoveProjectsState = { ...state };

      newRemoveProjectsState[action.payload.group] = _.omit(
        state[action.payload.group],
        action.payload.projects
      );

      const upstreamsToRemove = [];
      const downstreamsToRemove = [];

      // remove any upstreams or downstreams exclusive to the project from Group.UPSTREAM/Group.DOWNSTREAM
      if (action.payload.group == Group.PROJECTS) {
        action.payload.projects.forEach(p => {
          _.forIn(state.upstreamToProject, (v, k) => {
            // if an upstream is associated to the project
            if (v[p]) {
              // if the upstream is exclusive to project
              if (Object.keys(v).length == 1) {
                upstreamsToRemove.push(k);
                // remove as it was only tied to this project and b/c both the project and upstream will be removed from the state
                // otherwise we'll end up with stale information
                delete state.upstreamToProject[k];
              } else {
                // project is not exclusive to the upstream but let's remove as it will be removed from Group.PROJECTS state
                // otherwise we'll end up with stale information
                delete v[p];
              }
            }
          });
          _.forIn(state.downstreamToProject, (v, k) => {
            // if a downstream is associated to the project
            if (v[p]) {
              // if the downstream is exclusive to project
              if (Object.keys(v).length == 1) {
                downstreamsToRemove.push(k);
                // remove as it was only tied to this project and b/c both the project and downstream will be removed from the state
                // otherwise we'll end up with stale information
                delete state.downstreamToProject[k];
              } else {
                // project is not exclusive to the downstream but let's remove as it will be removed from Group.PROJECTS state
                // otherwise we'll end up with stale information
                delete v[p];
              }
            }
          });
        });

        if (upstreamsToRemove.length > 0) {
          newRemoveProjectsState[Group.UPSTREAM] = _.omit(state[Group.UPSTREAM], upstreamsToRemove);
        }
        if (downstreamsToRemove.length > 0) {
          newRemoveProjectsState[Group.DOWNSTREAM] = _.omit(
            state[Group.DOWNSTREAM],
            downstreamsToRemove
          );
        }
      }

      return newRemoveProjectsState;
    }
    case "TOGGLE_PROJECTS": {
      // TODO: hide exclusive upstreams and downstreams if group is PROJECTS
      return {
        ...state,
        [action.payload.group]: {
          ...state[action.payload.group],
          ...Object.fromEntries(
            action.payload.projects.map(key => [
              key,
              {
                ...state[action.payload.group][key],
                checked: !state[action.payload.group][key].checked,
              },
            ])
          ),
        },
      };
    }
    case "ONLY_PROJECTS": {
      const newState = { ...state };

      newState[action.payload.group] = Object.fromEntries(
        Object.keys(state[action.payload.group]).map(key => [
          key,
          { ...state[action.payload.group][key], checked: action.payload.projects.includes(key) },
        ])
      );

      return newState;
    }
    case "TOGGLE_ENTIRE_GROUP": {
      const newCheckedValue = !deriveSwitchStatus(state, action.payload.group);
      const newState = { ...state };
      newState[action.payload.group] = Object.fromEntries(
        Object.keys(state[action.payload.group]).map(key => [
          key,
          { ...state[action.payload.group][key], checked: newCheckedValue },
        ])
      );

      return newState;
    }
    // Background actions.

    case "HYDRATE_START": {
      return { ...state, loading: true };
    }
    case "HYDRATE_END": {
      let newState = { ...state, loading: false, error: undefined };

      _.forIn(
        action.payload.result as IClutch.project.v1.IGetProjectsResponse,
        (v: IClutch.project.v1.IProjectResult, k: string) => {
          // user owned project vs custom project
          if (v.from.users.length > 0) {
            newState = updateGroupstate(newState, Group.PROJECTS, k, { checked: true });
          } else if (v.from.selected) {
            newState = updateGroupstate(newState, Group.PROJECTS, k, {
              checked: true,
              custom: true,
            });
          }

          // add each upstream/downstream for the selected or user project
          if (v.from.users.length > 0 || v.from.selected) {
            v.project.dependencies.upstreams[
              "type.googleapis.com/clutch.core.project.v1.Project"
            ]?.ids.forEach(upstreamDep => {
              newState = updateGroupstate(newState, Group.UPSTREAM, upstreamDep, {
                checked: false,
              });
              if (!state.upstreamToProject[upstreamDep]) {
                state.upstreamToProject[upstreamDep] = { [k]: true };
              } else {
                state.upstreamToProject[upstreamDep][k] = true;
              }
            });

            v.project.dependencies.downstreams[
              "type.googleapis.com/clutch.core.project.v1.Project"
            ]?.ids.forEach(downstreamDep => {
              newState = updateGroupstate(newState, Group.DOWNSTREAM, downstreamDep, {
                checked: false,
              });
              if (!state.downstreamToProject[downstreamDep]) {
                state.downstreamToProject[downstreamDep] = { [k]: true };
              } else {
                state.downstreamToProject[downstreamDep][k] = true;
              }
            });
          }

          // stores the raw project data for each project in the API result
          newState.projectData[k] = v.project;
        }
      );
      return newState;
    }
    case "HYDRATE_ERROR":
      /*
       TODO: do we want to handle the error state differently? For example, when we render the error on the UI,
       it won't disapper unless there's a successful API call or if the user refreshes the page. If a user performs other
       actions, such as use the toggle/checkbox/ etc. the error message will be still be on the page

       TODO: when we add error handling for projects not found, we'll need to make sure we remove the not-found-project from project group
       (it's added automatically in the "ADD_PROJECTS" state)
      */
      return { ...state, loading: false, error: action.payload.result };
    default:
      throw new Error(`unknown resolver action`);
  }
};

export { selectorReducer };
