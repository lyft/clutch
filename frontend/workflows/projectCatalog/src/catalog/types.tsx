import type { clutch as IClutch } from "@clutch-sh/api";
import type { ClutchError } from "@clutch-sh/core";

export interface CatalogProject extends IClutch.core.project.v1.IProject {
  removable?: boolean;
}

type UserActionKind = "ADD_PROJECT" | "CLEAR_ERROR" | "REMOVE_PROJECT" | "SEARCH";

interface UserPayload {
  search?: string;
  projects?: IClutch.core.project.v1.IProject[];
}

interface UserAction {
  type: UserActionKind;
  payload?: UserPayload;
}

type BackgroundActionKind =
  | "HYDRATE_START"
  | "HYDRATE_END"
  | "HYDRATE_ERROR"
  | "SEARCH_START"
  | "SEARCH_END";

interface BackgroundPayload {
  result: any;
}

interface BackgroundAction {
  type: BackgroundActionKind;
  payload?: BackgroundPayload;
}

export type Action = UserAction | BackgroundAction;
export interface CatalogState {
  projects: CatalogProject[];
  search?: string;
  isLoading: boolean;
  isSearching: boolean;
  error: ClutchError | undefined;
}
