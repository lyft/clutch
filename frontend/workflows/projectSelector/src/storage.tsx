import _ from "lodash";

import type { GlobalProjectState, State } from "./types";
import { Group, isGlobalProjectState } from "./types";

const LOCAL_STORAGE_STATE_KEY = "dashState";

/**
 * Attempts to load state stored in local storage and merge it into the specified state.
 *
 * The specified state is returned unmodified if on of the following it met:
 *   * local storage does not have state
 *   * the stored state does not have the correct type
 *   * the stored state is not valid JSON
 */
const loadStoredState = (state: State): State => {
  // Grab stored state from local storage
  const storedState = window.localStorage.getItem(LOCAL_STORAGE_STATE_KEY);
  // If stored state does not exist return state unmodified
  if (!storedState) {
    return state;
  }

  try {
    const storedStateObject = JSON.parse(storedState);
    // If stored state is in the proper format merge it with existing state
    if (isGlobalProjectState(storedStateObject)) {
      // Merge will overwrite existing values in state with any found in the stored state
      return _.merge(state, storedStateObject);
    }
    // If stored state is not in the correct format purge it and return state unmodified
    window.localStorage.removeItem(LOCAL_STORAGE_STATE_KEY);
    return state;
  } catch {
    // If any errors occur return unmodified state
    return state;
  }
};

const storeState = (state: State) => {
  const localState = {
    [Group.PROJECTS]: state[Group.PROJECTS],
    [Group.UPSTREAM]: state[Group.UPSTREAM],
    [Group.DOWNSTREAM]: state[Group.DOWNSTREAM],
  } as GlobalProjectState;
  window.localStorage.setItem(LOCAL_STORAGE_STATE_KEY, JSON.stringify(localState));
};

export { loadStoredState, LOCAL_STORAGE_STATE_KEY, storeState };
