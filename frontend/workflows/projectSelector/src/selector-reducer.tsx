import type { clutch as IClutch } from "@clutch-sh/api";
import _ from "lodash";

import type { Action, State } from "./hello-world";
import { deriveSwitchStatus, Group, stateHelper } from "./hello-world";

const selectorReducer = (state: State, action: Action): State => {
  switch (action.type) {
    // User actions.

    case "ADD_PROJECTS": {
      const newAddProjectsState = { ...state };

      // a given custom project may already exist in the group so don't trigger a state update for those duplicates
      const uniqueCustomProjects = action.payload.projects.filter(
        (project: string) => !(project in newAddProjectsState[action.payload.group])
      );
      if (uniqueCustomProjects.length === 0) {
        return newAddProjectsState;
      }

      uniqueCustomProjects.forEach(v => {
        newAddProjectsState[Group.PROJECTS][v] = { checked: true, custom: true };
        // check if we already have project data for this project. if so, add the upstreamds/downstreams
        if (v in newAddProjectsState.projectData) {
          _.forIn(newAddProjectsState.projectData[v].dependencies.upstreams, v => {
            v.ids.forEach(v => {
              stateHelper.setChecked(newAddProjectsState, Group.UPSTREAM, v);
            });
          });

          _.forIn(newAddProjectsState.projectData[v].dependencies.downstreams, v => {
            v.ids.forEach(v => {
              stateHelper.setChecked(newAddProjectsState, Group.DOWNSTREAM, v);
            });
          });
        }
      });

      return newAddProjectsState;
    }
    case "REMOVE_PROJECTS": {
      // TODO: also remove any upstreams or downstreams related (only) to the project.
      // if group == Groups.PROJECT, hide exclusive downstream upstreams
      //
      return {
        ...state,
        [action.payload.group]: _.omit(state[action.payload.group], action.payload.projects),
      };
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
      const newState = { ...state, loading: false, error: undefined };

      _.forIn(
        action.payload.result as IClutch.project.v1.IGetProjectsResponse,
        (v: IClutch.project.v1.IProjectResult, k: string) => {
          // user owned project vs custom project
          if (v.from.users.length > 0) {
            stateHelper.setChecked(newState, Group.PROJECTS, k);
          } else if (v.from.selected) {
            stateHelper.setChecked(newState, Group.PROJECTS, k);
            newState[Group.PROJECTS][k].custom = true;
          }

          // add each upstream/downstream for the selected or user project
          if (v.from.users.length > 0 || v.from.selected) {
            _.forIn(v.project.dependencies.upstreams, v => {
              v.ids.forEach(v => {
                stateHelper.setChecked(newState, Group.UPSTREAM, v);
              });
            });
            _.forIn(v.project.dependencies.downstreams, v => {
              v.ids.forEach(v => {
                stateHelper.setChecked(newState, Group.DOWNSTREAM, v);
              });
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
