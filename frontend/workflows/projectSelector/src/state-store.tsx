import type { ClutchError } from "@clutch-sh/core";
import type { clutch as IClutch } from "@clutch-sh/api";
import type { GroupState, State } from "./hello-world"
import { Group } from "./hello-world"

export default class StateStore{
  [Group.PROJECTS]: GroupState;
  [Group.UPSTREAM]: GroupState;
  [Group.DOWNSTREAM]: GroupState;

  projectData: { [projectName: string]: IClutch.core.project.v1.IProject };
  loading: boolean;
  error: ClutchError | undefined;

  constructor(initialState: State){
    this[Group.PROJECTS] = initialState[Group.PROJECTS],
    this[Group.UPSTREAM] = initialState[Group.UPSTREAM],
    this[Group.DOWNSTREAM] = initialState[Group.DOWNSTREAM],
    this.projectData = initialState.projectData,
    this.loading = initialState.loading,
    this.error = initialState.error
  }
}

