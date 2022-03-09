import type { clutch as IClutch } from "@clutch-sh/api";

import type { HydratedData, StorageState } from "./storage-context";

type StorageActionKind = "STORE_DATA" | "REMOVE_DATA" | "EMPTY_DATA" | "HYDRATE";

type BackgroundStorageActionKind = "EMPTY_DATA" | "HYDRATE";

interface BackgroundPayload {
  data?: HydratedData;
}

interface UserPayload {
  componentName?: string;
  key?: string;
  data?: any;
  local?: boolean;
}

interface StorageAction {
  type: StorageActionKind;
  payload: UserPayload;
}

interface BackgroundAction {
  type: BackgroundStorageActionKind;
  payload?: BackgroundPayload;
}

type Action = StorageAction | BackgroundAction;

const rotateDataFromAPI = (data: IClutch.shortlink.v1.IShareableState[]): HydratedData => {
  const hydrated = {};

  data.forEach(({ key = "", state = {} }) => {
    hydrated[key] = state;
  });

  return hydrated;
};

const storeLocal = (key: string, data: any) => {
  try {
    window.localStorage.setItem(key, JSON.stringify(data));
  } catch (e) {
    // eslint-disable-next-line no-console
    console.error("Error saving to local storage", e);
  }
};

const removeLocal = (key: string) => {
  window.localStorage.removeItem(key);
};

const storageContextReducer = (state: StorageState, action: Action): StorageState => {
  switch (action.type) {
    case "STORE_DATA": {
      const { componentName, key, data, local = true } = action.payload as UserPayload;
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

      if (local) {
        storeLocal(key ?? componentName, data);
      }

      return { ...newState, ...tempStore };
    }
    case "REMOVE_DATA": {
      const { componentName, key, local = true } = action.payload as UserPayload;
      const newState = { ...state };
      const { tempStore } = newState;

      if (componentName && key) {
        delete tempStore[componentName][key];
      } else if (componentName) {
        delete tempStore[componentName];
      }

      if (local) {
        removeLocal(key ?? componentName);
      }

      return newState;
    }
    case "EMPTY_DATA": {
      return { ...state, tempStore: {} };
    }
    case "HYDRATE": {
      let newState = { ...state };

      if (action.payload && action.payload.data) {
        newState = {
          ...newState,
          store: rotateDataFromAPI(action.payload.data),
          shortLinked: true,
        };
      }

      return newState;
    }
    default:
      throw new Error("Unkown storage reducer action");
  }
};

export default storageContextReducer;
