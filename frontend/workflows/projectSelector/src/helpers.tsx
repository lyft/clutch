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

// evaluates whether an upstream/downstream should be hidden based on the checked status of/exclusivity to a project(s)
const hidden = (
  state: State,
  group: Group,
  project: string,
  dependency: string,
  depMapping: DependencyMappings
): boolean => {
  const hiddenDep = [];
  if (group === Group.UPSTREAM) {
    _.forIn(depMapping.upstreams, (v, k) => {
      if (v[project]) {
        // if project is unchecked and dependency is exclusive to that project, hide the dependency
        if (!state[Group.PROJECTS][project].checked && Object.keys(v).length === 1) {
          hiddenDep.push(k);
        } else if (Object.keys(v).every(p => !state[Group.PROJECTS][p].checked)) {
          // although the dependency is exclusive to more than 1 Group.Projects, if all Group.Projects share the dependency
          // are unchecked, it's safe to hide that dependency as well.
          hiddenDep.push(k);
        }
      }
    });
  } else if (group === Group.DOWNSTREAM) {
    _.forIn(depMapping.downstreams, (v, k) => {
      if (v[project]) {
        // if project is unchecked and dependency is exclusive to that project, hide the dependency
        if (!state[Group.PROJECTS][project].checked && Object.keys(v).length === 1) {
          hiddenDep.push(k);
        } else if (Object.keys(v).every(p => !state[Group.PROJECTS][p].checked)) {
          // although the dependency is exclusive to more than 1 Group.Projects, if all Group.Projects share the dependency
          // are unchecked, it's safe to hide that dependency as well.
          hiddenDep.push(k);
        }
      }
    });
  }

  if (hiddenDep.includes(dependency)) {
    return true;
  }

  return false;
};

/* filter out the state data on the following criteria:
if a project in Group.Projects is unchecked, remove the exclusive upstream(s)/downstream(s)
if all unchecked projects share a given upstream/downstream, remove the upstream/downstream
*/
const deriveStateData = (state: State): State => {
  // only go through the flow, if there are unchecked projects
  if (!Object.keys(state[Group.PROJECTS]).some(k => !state[Group.PROJECTS][k].checked)) {
    return state;
  }

  // get the relationships b/w upstreams/downstreams to projects
  const depMapping = dependencyToProjects(state, Group.PROJECTS);
  const hiddenUpstreams = [];
  const hiddenDownstreams = [];

  _.forEach(Object.keys(state[Group.PROJECTS]), project => {
    _.forEach(Object.keys(state[Group.UPSTREAM]), upstream => {
      if (hidden(state, Group.UPSTREAM, project, upstream, depMapping)) {
        hiddenUpstreams.push(upstream);
      }
    });

    _.forEach(Object.keys(state[Group.DOWNSTREAM]), downstream => {
      if (hidden(state, Group.DOWNSTREAM, project, downstream, depMapping)) {
        hiddenDownstreams.push(downstream);
      }
    });
  });

  const newState = { ...state };
  if (hiddenUpstreams.length > 0) {
    newState[Group.UPSTREAM] = _.omit(state[Group.UPSTREAM], hiddenUpstreams);
  }
  if (hiddenDownstreams.length > 0) {
    newState[Group.DOWNSTREAM] = _.omit(state[Group.DOWNSTREAM], hiddenDownstreams);
  }

  return newState;
};

export {
  deriveStateData,
  deriveSwitchStatus,
  DispatchContext,
  exclusiveProjectDependencies,
  PROJECT_TYPE_URL,
  StateContext,
  updateGroupstate,
  useDispatch,
  useReducerState,
};
