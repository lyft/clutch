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

type ComponentStorageActionKind = "STORE_DATA" | "REMOVE_DATA";

type HydrateStorageActionKind = "HYDRATE";

export interface HydratePayload {
  data?: IClutch.shortlink.v1.IShareableState[];
}

export interface ComponentPayload {
  componentName?: string;
  key?: string;
  data?: unknown;
  localStorage?: boolean;
}

interface ComponentStorageAction {
  type: ComponentStorageActionKind;
  payload: ComponentPayload;
}

interface HydrateStorageAction {
  type: HydrateStorageActionKind;
  payload?: HydratePayload;
}

export type Action = ComponentStorageAction | HydrateStorageAction;

export interface WorkflowStorageState {
  fromShortLink: boolean;
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
  fromShortLink: boolean;
  removeData: RemoveDataFn;
  retrieveData: RetrieveDataFn;
  storeData: StoreDataFn;
}

const defaultWorkflowStorageState: WorkflowStorageState = {
  fromShortLink: false,
  workflowStore: {},
  workflowSessionStore: {},
};

export { defaultWorkflowStorageState };
