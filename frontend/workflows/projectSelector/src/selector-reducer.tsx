import _ from "lodash";

import type { Action, State } from "./hello-world";
import { deriveSwitchStatus, Group } from "./hello-world";

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

      _.forIn(action.payload.result, (v, k) => {
        // TODO: handle v.from.user
        if (v.from.selected){
          state[Group.PROJECTS][k] = { checked: true, custom: true };
        }

        let upstreamsDeps = []
        // collect upstreams for each project in the results
        _.forIn(v.project.dependencies.upstreams, (v) => {
          upstreamsDeps = v.id
        });
        // TODO: handle v.from.user
        if (v.from.selected){
          upstreamsDeps.forEach(v => {
            // Add each upstream for the selected or user project
            state[Group.UPSTREAM][v] = { checked: false };
          })
        }

        let downstreamsDeps = []
        // collect downstreams for each project in the results
        _.forIn(v.project.dependencies.downstreams, (v) => {
          downstreamsDeps = v.id
        });
        // TODO: handle v.from.user
        if (v.from?.selected){
          downstreamsDeps.forEach(v => {
            // Add each downstream for the selected or user project
            state[Group.DOWNSTREAM][v] = { checked: false };
          })
        }

        // stores the raw project data for each project in the API result
        state.projectData[k] = {
          name: v.project.name,
          tier: v.project.tier,
          owners: v.project.owners,
          languages: v.project.languages,
          data: v.project.data,
          upstreams: upstreamsDeps,
          downstreams: downstreamsDeps,
        };



      });
      return newPostAPICallState;
    }
    // TODO: test out an error case
    case "HYDRATE_ERROR":
      return { ...state, loading: false, error: action.payload.result };
    default:
      throw new Error(`unknown resolver action`);
  }
};

export { selectorReducer };
