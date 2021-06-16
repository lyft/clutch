import _ from "lodash";

import type { Action, State } from "./hello-world";
import { deriveSwitchStatus, Group } from "./hello-world";

const selectorReducer = (state: State, action: Action): State => {
  switch (action.type) {
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
      const newOnlyProjectState = { ...state };

      newOnlyProjectState[action.payload.group] = Object.fromEntries(
        Object.keys(state[action.payload.group]).map(key => [
          key,
          { ...state[action.payload.group][key], checked: action.payload.projects.includes(key) },
        ])
      );

      return newOnlyProjectState;
    }
    case "TOGGLE_ENTIRE_GROUP": {
      const newCheckedValue = !deriveSwitchStatus(state, action.payload.group);
      const newGroupToggledState = { ...state };
      newGroupToggledState[action.payload.group] = Object.fromEntries(
        Object.keys(state[action.payload.group]).map(key => [
          key,
          { ...state[action.payload.group][key], checked: newCheckedValue },
        ])
      );

      return newGroupToggledState;
    }
    // Background actions.
    case "HYDRATE_START": {
      return { ...state, loading: true };
    }
    case "HYDRATE_END": {
      const newPostAPICallState = { ...state, loading: false };
      // TODO: handle payload.
      _.forIn(action.payload.result, (v, k) => {
        // Add each project to the projects list.
        state[Group.PROJECTS][k] = { checked: true };
        state.projectData[k] = {};

        // Add each upstream.
        v.upstreams.forEach(v => {
          state[Group.UPSTREAM][v] = { checked: false };
          state.projectData[v] = {};
        });

        // Add each downstream.
        v.downstreams.forEach(v => {
          state[Group.DOWNSTREAM][v] = { checked: false };
          state.projectData[v] = {};
        });

        // Update project data for each.
      });
      return newPostAPICallState;
    }
    default:
      throw new Error(`unknown resolver action`);
  }
};

export { selectorReducer };
