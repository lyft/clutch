import { removeLocalData, rotateDataFromAPI, storeLocalData } from "./helpers";
import type { Action, BackgroundPayload, UserPayload, WorkflowStorageState } from "./types";

/**
 * Reducer for the WorkflowStorageContext
 * This will act on the WorkflowStorageState and add / remove items from the temporary storage
 * as well as localStorage, this will (optionally) keep all storage actions in one location for all
 * workflows and allow for easier state hydration
 */
const workflowStorageContextReducer = (
  state: WorkflowStorageState,
  action: Action
): WorkflowStorageState => {
  switch (action.type) {
    // Will add data to our temporary storage as well as the local storage
    case "STORE_DATA": {
      const { componentName, key, data, localStorage = true } = action.payload as UserPayload;
      const newState = { ...state };
      const { shortLinked, tempStore } = newState;

      if (!tempStore[componentName]) {
        tempStore[componentName] = {};
      }

      if (key.length) {
        tempStore[componentName][key] = data;
      } else {
        tempStore[componentName] = { ...tempStore[componentName], ...(data as any) };
      }

      if (localStorage && !shortLinked) {
        storeLocalData(key ?? componentName, data);
      }

      return { ...newState, ...tempStore };
    }
    // Will remove data from our temporary storage as well as the local storage
    case "REMOVE_DATA": {
      const { componentName, key, localStorage = true } = action.payload as UserPayload;
      const newState = { ...state };
      const { shortLinked, tempStore } = newState;

      if (componentName && key) {
        delete tempStore[componentName][key];
      } else if (componentName) {
        delete tempStore[componentName];
      }

      if (localStorage && !shortLinked) {
        removeLocalData(key ?? componentName);
      }

      return newState;
    }
    // Will take a given input of data from an API and add it to the state as 'store', the only time this portion of the state should ever be modified
    case "HYDRATE": {
      const { data } = action.payload as BackgroundPayload;

      if (data) {
        return {
          ...state,
          store: rotateDataFromAPI(data),
          shortLinked: true,
        };
      }

      return state;
    }
    default:
      throw new Error("Unknown workflow storage reducer action");
  }
};

export default workflowStorageContextReducer;
