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

type BackgroundStorageActionKind = "HYDRATE";

export interface BackgroundPayload {
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

export interface WorkflowStorageState {
  shortLinked: boolean;
  store: HydratedData;
  tempStore: HydratedData;
}

export interface WorkflowStorageContextProps {
  shortLinked: boolean;
  removeData: (componentName: string, key: string, localStorage?: boolean) => void;
  retrieveData: (componentName: string, key: string, defaultData?: any) => any;
  storeData: (componentName: string, key: string, data: any, localStorage?: boolean) => void;
}

const defaultWorkflowStorageState: WorkflowStorageState = {
  shortLinked: false,
  store: {},
  tempStore: {},
};

export { defaultWorkflowStorageState };
