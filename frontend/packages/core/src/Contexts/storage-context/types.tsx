import type { clutch as IClutch } from "@clutch-sh/api";

export interface HydrateData {
  route: string;
  data: HydratedData;
}

// data is meant to be stored in a manner of:
/**
 * {
 *      componentName1: {
 *          key1: {
 *              data: ...
 *          },
 *          key2: {
 *              data: ...
 *          },
 *      }
 * }
 */
export interface HydratedData {
  [key: string]: {
    [key: string]: any;
  };
}

type StorageActionKind = "STORE_DATA" | "REMOVE_DATA";

type BackgroundStorageActionKind = "EMPTY_TEMP_DATA" | "HYDRATE";

interface BackgroundPayload {
  data?: IClutch.shortlink.v1.IShareableState[];
}

export interface UserPayload {
  componentName?: string;
  key?: string;
  data?: any;
  localStorage?: boolean;
}

interface StorageAction {
  type: StorageActionKind;
  payload: UserPayload;
}

interface BackgroundAction {
  type: BackgroundStorageActionKind;
  payload?: BackgroundPayload;
}

export type Action = StorageAction | BackgroundAction;

export interface StorageState {
  shortLinked: boolean;
  store: HydratedData;
  tempStore: HydratedData;
}

export type StoreDataFn = (
  componentName: string,
  key: string,
  data: any,
  localStorage?: boolean
) => void;
export type StoreLocalDataFn = (key: string, data: any) => void;
export type RemoveDataFn = (componentName: string, key: string, localStorage?: boolean) => void;
export type RemoveLocalDataFn = (key: string) => void;
export type RetrieveDataFn = (componentName: string, key: string, defaultData?: any) => any;
export type RetrieveLocalDataFn = (key: string, defaultData?: any) => any;
export type ClearTempDataFn = () => void;
export type TempDataFn = () => HydratedData;

export interface StorageContextProps {
  shortLinked: boolean;
  functions: {
    storeData: StoreDataFn;
    storeLocalData: StoreLocalDataFn;
    removeData: RemoveDataFn;
    removeLocalData: RemoveLocalDataFn;
    retrieveData: RetrieveDataFn;
    retrieveLocalData: RetrieveLocalDataFn;
    clearTempData: ClearTempDataFn;
    tempData: TempDataFn;
  };
}

const defaultStorageState: StorageState = {
  shortLinked: false,
  store: {},
  tempStore: {},
};

export { defaultStorageState };
