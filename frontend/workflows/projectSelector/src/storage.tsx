import _ from "lodash";

import type { GlobalProjectState, State } from "./types";
import { Group, isGlobalProjectState } from "./types";

export const COMPONENT_NAME = "ProjectSelector";
const STORAGE_STATE_KEY = "dashState";

/**
 * Attempts to load state stored in local storage and merge it into the specified state.
 *
 * The specified state is returned unmodified if on of the following it met:
 *   * local storage does not have state
 *   * the stored state does not have the correct type
 *   * the stored state is not valid JSON
 */
const loadStoredState = (retrieve, remove, state: State): State => {
  // Grab stored state
  const storedState = retrieve(COMPONENT_NAME, STORAGE_STATE_KEY);
  // Grab stored state from local storage
  // const storedState = window.localStorage.getItem(STORAGE_STATE_KEY);
  // If stored state does not exist return state unmodified
  if (!storedState) {
    return state;
  }

  try {
    // If stored state is in the proper format merge it with existing state
    if (isGlobalProjectState(storedState)) {
      // Merge will overwrite existing values in state with any found in the stored state
      return _.merge(state, storedState);
    }
    // If stored state is not in the correct format purge it and return state unmodified
    remove(COMPONENT_NAME, STORAGE_STATE_KEY);
    return state;
  } catch {
    // If any errors occur return unmodified state
    return state;
  }
};

const storeState = (store, state: State, local: boolean) => {
  const localState = {
    [Group.PROJECTS]: state[Group.PROJECTS],
    [Group.UPSTREAM]: state[Group.UPSTREAM],
    [Group.DOWNSTREAM]: state[Group.DOWNSTREAM],
  } as GlobalProjectState;
  store(COMPONENT_NAME, STORAGE_STATE_KEY, localState, local);
};

export { loadStoredState, storeState };
