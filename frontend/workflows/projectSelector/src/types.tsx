import type { clutch as IClutch } from "@clutch-sh/api";
import type { ClutchError } from "@clutch-sh/core";

export enum Group {
  PROJECTS,
  UPSTREAM,
  DOWNSTREAM,
}

type UserActionKind =
  | "ADD_PROJECTS"
  | "REMOVE_PROJECTS"
  | "TOGGLE_PROJECTS"
  | "TOGGLE_ENTIRE_GROUP"
  | "ONLY_PROJECTS";

interface UserAction {
  type: UserActionKind;
  payload: UserPayload;
}

interface UserPayload {
  group: Group;
  projects?: string[];
}

type BackgroundActionKind = "HYDRATE_START" | "HYDRATE_END" | "HYDRATE_ERROR";

interface BackgroundAction {
  type: BackgroundActionKind;
  payload?: BackgroundPayload;
}

interface BackgroundPayload {
  result: any;
}

export type Action = BackgroundAction | UserAction;

export interface ProjectsState {
  [Group.PROJECTS]: GroupState;
  [Group.UPSTREAM]: GroupState;
  [Group.DOWNSTREAM]: GroupState;
}

/**
 * Determines if an object is of type @type {ProjectsState}.
 */
 const isProjectsState = (state: ProjectsState | Object): state is ProjectsState => {
  const projectsState = state as ProjectsState;
  const projects = projectsState[Group.PROJECTS];
  const upstream = projectsState[Group.UPSTREAM];
  const downstream = projectsState[Group.DOWNSTREAM];

  const hasProjects = isGroupState(projects);
  const hasUpstream = isGroupState(upstream);
  const hasDownstream = isGroupState(downstream);
  return hasProjects && hasUpstream && hasDownstream;
};

export interface State extends ProjectsState {
  projectData: { [projectName: string]: IClutch.core.project.v1.IProject };
  loading: boolean;
  error: ClutchError | undefined;
}

export interface GroupState {
  [projectName: string]: ProjectState;
}

/**
 * Determines if an object is of type @type {GroupState}.
 */
 const isGroupState = (state: GroupState | object | undefined): state is GroupState => {
  if (state === undefined) {
    return false;
  }
  const projectStates = Object.values(state as GroupState);
  return projectStates.filter(s => isProjectState(s)).length === projectStates.length;
};


// n.b. if you are updating ProjectState be sure to update the custom type guard below it.
export interface ProjectState {
  checked: boolean;
  // TODO: hidden should be derived?
  hidden?: boolean; // upstreams and downstreams are hidden when their parent is unchecked unless other parents also use them.
  custom?: boolean;
}

/**
 * Determines if an object is of type @type {ProjectState}.
 */
 const isProjectState = (state: ProjectState | object): state is ProjectState => {
  const pState = (state as ProjectState);
  const checkedProp = pState?.checked;
  const hasRequiredProps = checkedProp !== undefined && typeof checkedProp === "boolean";
  
  const validOptionalProps = (
    (pState?.hidden !== undefined ? typeof pState.hidden === "boolean" : true) &&
    (pState?.custom !== undefined ? typeof pState.custom === "boolean" : true)
  );
  return hasRequiredProps && validOptionalProps;
};

export interface DashState {
  // Contains the names of selected projects, upstreams, and downstreams merged together.
  selected: string[];

  // Contains a map of project names to the full project data.
  projectData: { [projectName: string]: IClutch.core.project.v1.IProject };
}

export type DashActionKind = "UPDATE_SELECTED";

export interface DashAction {
  type: DashActionKind;
  payload: DashState;
}

export {isGroupState, isProjectState, isProjectsState };