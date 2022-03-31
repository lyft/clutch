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
    [key: string]: unknown;
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
  data?: unknown;
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
  workflowStore: HydratedData;
  workflowSessionStore: HydratedData;
}

export type RemoveDataFn = (componentName: string, key: string, localStorage?: boolean) => void;
export type RetrieveDataFn = (componentName: string, key: string, defaultData?: unknown) => unknown;
export type StoreDataFn = (
  componentName: string,
  key: string,
  data: unknown,
  localStorage?: boolean
) => void;

export interface WorkflowStorageContextProps {
  shortLinked: boolean;
  removeData: RemoveDataFn;
  retrieveData: RetrieveDataFn;
  storeData: StoreDataFn;
}

const defaultWorkflowStorageState: WorkflowStorageState = {
  shortLinked: false,
  workflowStore: {},
  workflowSessionStore: {},
};

export { defaultWorkflowStorageState };
