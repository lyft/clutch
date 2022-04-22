import type { clutch as IClutch } from "@clutch-sh/api";

import type { HydratedData } from "./types";

// API data comes in as an array, this will rotate it into an object usable by the StorageContext
const transformAPISharedState = (data: IClutch.shortlink.v1.IShareableState[]): HydratedData => {
  const hydrated: HydratedData = {};

  data.forEach(({ key, state = {} }) => {
    hydrated[key] = state as any;
  });

  return hydrated;
};

const storeLocalData = (key: string, data: unknown) => {
  try {
    window.localStorage.setItem(key, JSON.stringify(data));
  } catch (e) {
    // eslint-disable-next-line no-console
    console.error("Error saving to local storage", e);
  }
};

const removeLocalData = (key: string) => window.localStorage.removeItem(key);

const retrieveLocalData = (key: string) => {
  const localData = window.localStorage.getItem(key);

  if (localData) {
    try {
      return JSON.parse(localData);
    } catch (_) {
      return localData;
    }
  }
};

type GenericRetrieve = <T>(
  store: HydratedData,
  componentName: string,
  key: string,
  defaultData: T
) => T;

const retrieveData: GenericRetrieve = (
  store: HydratedData,
  componentName: string,
  key: string,
  defaultData?
) => {
  if (store && store[componentName]) {
    return key.length ? store[componentName][key] : store[componentName];
  }

  if (key.length) {
    return retrieveLocalData(key) ?? defaultData;
  }

  return defaultData;
};

export { removeLocalData, retrieveData, storeLocalData, transformAPISharedState };
