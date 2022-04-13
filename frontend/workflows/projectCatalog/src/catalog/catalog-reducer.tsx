import type { Action, CatalogState } from "./types";

// TODO: Migrate all local storage actions from ./storage.tsx to here
const catalogReducer = (state: CatalogState, action: Action): CatalogState => {
  switch (action.type) {
    case "ADD_PROJECT": {
      return { ...state, projects: action?.payload?.projects || [] };
    }
    case "REMOVE_PROJECT": {
      return { ...state, projects: action?.payload?.projects || [] };
    }
    case "SEARCH": {
      return { ...state, search: action?.payload?.search };
    }
    case "SEARCH_START": {
      return { ...state, error: undefined, isSearching: true };
    }
    case "SEARCH_END": {
      return { ...state, search: "", isSearching: false };
    }
    case "HYDRATE_START": {
      return { ...state, isLoading: true };
    }
    case "HYDRATE_END": {
      return { ...state, projects: action?.payload?.result, isLoading: false };
    }
    case "HYDRATE_ERROR": {
      return { ...state, error: action?.payload?.result, isLoading: false };
    }
    default:
      throw new Error("Unknown action");
  }
};

export default catalogReducer;
