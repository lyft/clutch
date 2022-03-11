import { removeLocalData, rotateDataFromAPI, storeLocalData } from "./helpers";
import type { Action, BackgroundPayload, StorageState, UserPayload } from "./types";

/**
 * Reducer for the StorageContext
 * This will act on the StorageState and add / remove items from the temporary storage
 * as well as localStorage, this will keep all storage actions in one location for all
 * components and allow for easier state hydration across the entire application
 */
const storageContextReducer = (state: StorageState, action: Action): StorageState => {
  switch (action.type) {
    // Will add data to our temporary storage as well as the local storage
    case "STORE_DATA": {
      const { componentName, key, data, localStorage = true } = action.payload as UserPayload;
      const newState = { ...state };
      const { tempStore } = newState;

      if (!tempStore[componentName]) {
        tempStore[componentName] = {};
      }

      if (key.length) {
        tempStore[componentName][key] = data;
      } else {
        tempStore[componentName] = { ...tempStore[componentName], ...data };
      }

      if (localStorage) {
        storeLocalData(key ?? componentName, data);
      }

      return { ...newState, ...tempStore };
    }
    // Will remove data from our temporary storage as well as the local storage
    case "REMOVE_DATA": {
      const { componentName, key, localStorage = true } = action.payload as UserPayload;
      const newState = { ...state };
      const { tempStore } = newState;

      if (componentName && key) {
        delete tempStore[componentName][key];
      } else if (componentName) {
        delete tempStore[componentName];
      }

      if (localStorage) {
        removeLocalData(key ?? componentName);
      }

      return newState;
    }
    // Will take an input route and check to see if we're already shortlinked and if we are, will only clear if the route is different
    case "CLEAR_SHORT_LINK": {
      const { route = "" } = action.payload as BackgroundPayload;

      if (state.shortLinked && state.shortLinkedRoute !== route) {
        return { ...state, store: {}, shortLinked: false, shortLinkedRoute: undefined };
      }

      return state;
    }
    // Will clear out the temporary storage
    case "EMPTY_TEMP_DATA": {
      return { ...state, tempStore: {} };
    }
    // Will take a given input of data from an API and add it to the state as 'store', the only time this portion of the state should ever be modified
    case "HYDRATE": {
      const { data, route } = action.payload as BackgroundPayload;

      if (data && route) {
        return {
          ...state,
          store: rotateDataFromAPI(data),
          shortLinked: true,
          shortLinkedRoute: route,
        };
      }

      return state;
    }
    default:
      throw new Error("Unknown storage reducer action");
  }
};

export default storageContextReducer;
