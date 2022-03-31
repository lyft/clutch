import _ from "lodash";

import type { GlobalProjectState, State } from "./types";
import { Group, isGlobalProjectState } from "./types";

export const STORAGE_STATE_KEY = "dashState";
export const COMPONENT_NAME = "ProjectSelector";

/**
 * Attempts to load state stored in local storage and merge it into the specified state.
 *
 * The specified state is returned unmodified if on of the following it met:
 *   * local storage does not have state
 *   * the stored state does not have the correct type
 *   * the stored state is not valid JSON
 */
const loadStoredState = (state: State, retrieveData, removeData): State => {
  // Grab stored state from local storage
  const storedState = retrieveData(COMPONENT_NAME, STORAGE_STATE_KEY, true);
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
    removeData(COMPONENT_NAME, STORAGE_STATE_KEY, true);
    return state;
  } catch {
    // If any errors occur return unmodified state
    return state;
  }
};

const getLocalState = (state: State): GlobalProjectState => {
  const localState = {
    [Group.PROJECTS]: state[Group.PROJECTS],
    [Group.UPSTREAM]: state[Group.UPSTREAM],
    [Group.DOWNSTREAM]: state[Group.DOWNSTREAM],
  } as GlobalProjectState;

  return localState;
};

export { loadStoredState, getLocalState };
