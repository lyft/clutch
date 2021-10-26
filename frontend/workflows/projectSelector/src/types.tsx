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

export interface GlobalProjectState {
  [Group.PROJECTS]: GroupState;
  [Group.UPSTREAM]: GroupState;
  [Group.DOWNSTREAM]: GroupState;
}

export interface GroupState {
  [projectName: string]: ProjectState;
}

// n.b. if you are updating ProjectState be sure to update the custom type guard below it.
export interface ProjectState {
  checked: boolean;
  custom?: boolean;
}

/**
 * Determines if an object is of type @type {ProjectState}.
 */
const isProjectState = (state: ProjectState | object): state is ProjectState => {
  const pState = state as ProjectState;
  const checkedProp = pState?.checked;
  const hasRequiredProps = checkedProp !== undefined && typeof checkedProp === "boolean";

  const validOptionalProps =
    pState?.custom !== undefined ? typeof pState.custom === "boolean" : true;
  return hasRequiredProps && validOptionalProps;
};

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

/**
 * Determines if an object is of type @type {GlobalState}.
 */
const isGlobalProjectState = (state: GlobalProjectState | Object): state is GlobalProjectState => {
  const globalState = state as GlobalProjectState;
  const projects = globalState[Group.PROJECTS];
  const upstream = globalState[Group.UPSTREAM];
  const downstream = globalState[Group.DOWNSTREAM];

  const hasProjects = isGroupState(projects);
  const hasUpstream = isGroupState(upstream);
  const hasDownstream = isGroupState(downstream);
  return hasProjects && hasUpstream && hasDownstream;
};

export interface State extends GlobalProjectState {
  projectData: { [projectName: string]: IClutch.core.project.v1.IProject };
  loading: boolean;
  error: ClutchError | undefined;
}

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

export interface TimeData {
  /** eventsKey corresponds to entity owning the data - i.e. card name */
  [eventsKey: string]: EventData;
}

export interface EventData {
  /**
   * Mapping of project names to their event time points
   * See https://github.com/lyft/clutch/blob/main/api/timeseries/v1/timeseries.proto
   */
  points: { [projectName: string]: IClutch.timeseries.v1.IPoint[] };
  /** The emoji that will be used for this card on the event timeline */
  emoji: string;
}

/** Used by the reducer to update the time data in our context. */
export interface TimeDataUpdate {
  /** The name of the card or entity that is updating */
  key: string;
  /** The projects with their timeseries data and emoji */
  eventData: EventData;
}

export interface TimelineState {
  timeData: TimeData;
}

export type TimelineActionKindUpdate = "UPDATE";
export interface TimelineAction {
  type: TimelineActionKindUpdate;
  payload: TimeDataUpdate;
}
export interface TimeRangeState {
  /** The start of the time window the user has chosen in milliseconds */
  startTimeMs: number;
  /** The end of time window the user has chosen in milliseconds */
  endTimeMs: number;
}

export type TimeRangeActionUpdate = "UPDATE";

export interface TimeRangeAction {
  type: TimeRangeActionUpdate;
  payload: TimeRangeState;
}

export { isGroupState, isProjectState, isGlobalProjectState };
