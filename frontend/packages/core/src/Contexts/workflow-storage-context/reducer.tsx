import { removeLocalData, storeLocalData, transformAPISharedState } from "./helpers";
import type { Action, ComponentPayload, HydratePayload, WorkflowStorageState } from "./types";

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
      const { componentName, key, data, localStorage = true } = action.payload as ComponentPayload;
      const newState = { ...state };
      const { fromShortLink, workflowSessionStore } = newState;

      if (!componentName || !componentName.length) {
        return state;
      }

      if (!workflowSessionStore[componentName]) {
        workflowSessionStore[componentName] = {};
      }

      if (key.length) {
        workflowSessionStore[componentName][key] = data;
      } else {
        workflowSessionStore[componentName] = {
          ...workflowSessionStore[componentName],
          ...(data as any),
        };
      }

      if (localStorage && !fromShortLink) {
        storeLocalData(key ?? componentName, data);
      }

      return { ...newState, workflowSessionStore };
    }
    // Will remove data from our temporary storage as well as the local storage
    case "REMOVE_DATA": {
      const { componentName, key, localStorage = true } = action.payload as ComponentPayload;
      const newState = { ...state };
      const { fromShortLink, workflowSessionStore } = newState;

      if (!componentName || !componentName.length) {
        return state;
      }

      if (componentName && key) {
        delete workflowSessionStore[componentName][key];
      } else if (componentName) {
        delete workflowSessionStore[componentName];
      }

      if (localStorage && !fromShortLink) {
        removeLocalData(key ?? componentName);
      }

      return newState;
    }
    // Will take a given input of data from an API and add it to the state as 'store', the only time this portion of the state should ever be modified
    case "HYDRATE": {
      const { data } = action.payload as HydratePayload;

      if (data) {
        return {
          ...state,
          fromShortLink: true,
          workflowStore: transformAPISharedState(data),
        };
      }

      return state;
    }
    default:
      throw new Error("Unknown workflow storage reducer action");
  }
};

export default workflowStorageContextReducer;
