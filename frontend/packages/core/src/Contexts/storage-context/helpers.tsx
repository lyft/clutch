import type { clutch as IClutch } from "@clutch-sh/api";

import type { HydratedData, StorageState } from "./types";

// API data comes in as an array, this will rotate it into an object usable by the StorageContext
const rotateDataFromAPI = (data: IClutch.shortlink.v1.IShareableState[]): HydratedData => {
  const hydrated: HydratedData = {};

  data.forEach(({ key = "", state = {} }) => {
    hydrated[key] = state;
  });

  return hydrated;
};

const storeLocalData = (key: string, data: any) => {
  try {
    window.localStorage.setItem(key, JSON.stringify(data));
  } catch (e) {
    // eslint-disable-next-line no-console
    console.error("Error saving to local storage", e);
  }
};

const removeLocalData = (key: string) => window.localStorage.removeItem(key);

const retrieveLocalData = (key: string) => window.localStorage.getItem(key);

const retrieveData = (
  storageState: StorageState,
  componentName: string,
  key: string,
  defaultData?: any
): any => {
  const { store } = storageState;

  if (store && store[componentName]) {
    return key.length ? store[componentName][key] : store[componentName];
  }

  if (key.length) {
    const localData = retrieveLocalData(key);

    if (localData) {
      try {
        return JSON.parse(localData);
      } catch (_) {
        return localData;
      }
    }
  }

  return defaultData;
};

export { rotateDataFromAPI, removeLocalData, retrieveData, retrieveLocalData, storeLocalData };
