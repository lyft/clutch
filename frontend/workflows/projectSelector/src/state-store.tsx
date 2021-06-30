import type { clutch as IClutch } from "@clutch-sh/api";
import type { ClutchError } from "@clutch-sh/core";

import type { GroupState, State } from "./hello-world";
import { Group } from "./hello-world";

export default class StateStore implements State {
  [Group.PROJECTS]: GroupState;

  [Group.UPSTREAM]: GroupState;

  [Group.DOWNSTREAM]: GroupState;

  projectData: { [projectName: string]: IClutch.core.project.v1.IProject };

  loading: boolean;

  error: ClutchError | undefined;

  constructor() {
    this[Group.PROJECTS] = {};
    this[Group.UPSTREAM] = {};
    this[Group.DOWNSTREAM] = {};
    this.projectData = {};
    this.loading = false;
    this.error = undefined;
  }

  setChecked(state: State, group: Group, project: string) {
    // no matter the group, we preserve the checked state
    if (project in state[group]) {
      state[group][project].checked = state[group][project].checked;
    } else if (group == Group.PROJECTS) {
      // projects in this group are checked true by default
      state[group][project] = { checked: true };
    } else {
      // projects in the upstream/downstream groups are checked false by default
      state[group][project] = { checked: false };
    }
  }

}
