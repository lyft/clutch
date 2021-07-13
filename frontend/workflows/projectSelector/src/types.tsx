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
  | "ONLY_PROJECTS"
  | "HIDE_PROJECTS";

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

export interface State {
  [Group.PROJECTS]: GroupState;
  [Group.UPSTREAM]: GroupState;
  [Group.DOWNSTREAM]: GroupState;

  projectData: { [projectName: string]: IClutch.core.project.v1.IProject };
  loading: boolean;
  error: ClutchError | undefined;
}

interface GroupState {
  [projectName: string]: ProjectState;
}

// TODO: add new field to designate whether a Group.Project has hidden upstreams/downstreams so that way we can
// support hiding all upstreams/downstreams if they are shared by other Group.Project who have already opted to
// hide their dependencies
export interface ProjectState {
  checked: boolean;
  // TODO: hidden should be derived?
  hidden?: boolean; // upstreams and downstreams are hidden when their parent is unchecked unless other parents also use them.
  custom?: boolean;
}
