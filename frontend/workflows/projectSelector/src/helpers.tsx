import * as React from "react";
import _ from "lodash";

import type { Action, ProjectState, State } from "./types";
import { Group } from "./types";

const PROJECT_TYPE_URL = "type.googleapis.com/clutch.core.project.v1.Project";

interface DependencyMappings {
  upstreams?: { [dependency: string]: { [project: string]: boolean } };
  downstreams?: { [dependency: string]: { [project: string]: boolean } };
}

const StateContext = React.createContext<State | undefined>(undefined);
const useReducerState = () => {
  return React.useContext(StateContext);
};

const DispatchContext = React.createContext<(action: Action) => void | undefined>(undefined);
const useDispatch = () => {
  return React.useContext(DispatchContext);
};

// TODO(perf): call with useMemo().
const deriveSwitchStatus = (state: State, group: Group): boolean => {
  return (
    Object.keys(state[group]).length > 0 &&
    Object.keys(state[group]).every(key => state[group][key].checked)
  );
};

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

const updateHiddenState = (state: State, group: Group, project: string, hidden: boolean): State => {
  const newState = { ...state };
  newState[group][project].hidden = hidden;
  return newState;
};

const dependencyToProjects = (state: State, group: Group): DependencyMappings => {
  const upstreamMap = {};
  const downstreamMap = {};

  const projects = Object.keys(state[group]);
  projects.forEach(p => {
    const { upstreams, downstreams } = state.projectData[p]?.dependencies || {};
    upstreams?.[PROJECT_TYPE_URL]?.ids.forEach(u => {
      if (!upstreamMap[u]) {
        upstreamMap[u] = { [p]: true };
      } else {
        upstreamMap[u][p] = true;
      }
    });
    downstreams?.[PROJECT_TYPE_URL]?.ids.forEach(d => {
      if (!downstreamMap[d]) {
        downstreamMap[d] = { [p]: true };
      } else {
        downstreamMap[d][p] = true;
      }
    });
  });

  return { upstreams: upstreamMap, downstreams: downstreamMap };
};

const exclusiveProjectDependencies = (
  state: State,
  group: Group,
  project: string
): { upstreams: string[]; downstreams: string[] } => {
  const dependencyMap = dependencyToProjects(state, group);

  const upstreams = [];
  const downstreams = [];
  _.forIn(dependencyMap.upstreams, (v, k) => {
    if (v[project] && Object.keys(v).length === 1) {
      upstreams.push(k);
    }
  });

  _.forIn(dependencyMap.downstreams, (v, k) => {
    if (v[project] && Object.keys(v).length === 1) {
      downstreams.push(k);
    }
  });
  return { upstreams, downstreams };
};

// deriveHiddenStatus evaluates whether to hide an upstream/downstream based on the checked status of/exclusivity to projects in Group.Projects
const deriveHiddenStatus = (state: State, group: Group, key: string): boolean => {
  // we only support hidding an upstream/downstream
  if (group === Group.PROJECTS) {
    return false;
  }

  // only go through the rest of the flow, if there are unchecked projects
  const unCheckedProjects = Object.keys(state[Group.PROJECTS]).filter(
    k => !state[Group.PROJECTS][k].checked
  );
  if (unCheckedProjects.length === 0) {
    updateHiddenState(state, group, key, false);
    return false;
  }

  // get the relationships b/w upstreams/downstreams to project(s)
  const dependencyMapping = dependencyToProjects(state, Group.PROJECTS);

  const exclusive = [];
  _.forIn(state[Group.PROJECTS], (groupState, project) => {
    if (group === Group.UPSTREAM) {
      _.forIn(dependencyMapping.upstreams, (v, k) => {
        if (v[project]) {
          // if project is unchecked and dependency is exclusive to that project
          if (!groupState.checked && Object.keys(v).length === 1) {
            exclusive.push(k);
          } else if (Object.keys(v).every(p => !state[Group.PROJECTS][p].checked)) {
            // although the dependency is exclusive to more than 1 Group.Projects, if all Group.Projects share the dependency
            // are unchecked, it's safe to hide that dependency as well.
            exclusive.push(k);
          }
        }
      });
    } else if (group === Group.DOWNSTREAM) {
      _.forIn(dependencyMapping.downstreams, (v, k) => {
        if (v[project]) {
          // if project is unchecked and dependency is exclusive to that project
          if (!groupState.checked && Object.keys(v).length === 1) {
            exclusive.push(k);
          } else if (Object.keys(v).every(p => !state[Group.PROJECTS][p].checked)) {
            // although the dependency is exclusive to more than 1 Group.Projects, if all Group.Projects share the dependency
            // are unchecked, it's safe to hide that dependency as well.
            exclusive.push(k);
          }
        }
      });
    }
  });

  if (exclusive.includes(key)) {
    updateHiddenState(state, group, key, true);
    return true;
  }

  updateHiddenState(state, group, key, false);
  return false;
};

export {
  deriveHiddenStatus,
  deriveSwitchStatus,
  DispatchContext,
  exclusiveProjectDependencies,
  PROJECT_TYPE_URL,
  StateContext,
  updateGroupstate,
  useDispatch,
  useReducerState,
};
