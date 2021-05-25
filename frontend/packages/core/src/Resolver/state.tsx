import React from "react";
import type { clutch } from "@clutch-sh/api";

import type { ClutchError } from "../Network/errors";

enum ResolverAction {
  SCHEMAS_LOADING,
  SCHEMAS_SUCCCESS,
  SCHEMAS_ERROR,
  RESOLVING,
  RESOLVE_ERROR,
  RESOLVE_SUCCESS,
}

const initialState = {
  schemasLoading: true,
  allSchemas: [],
  searchableSchemas: [],
  schemaFetchError: undefined,
  resolverLoading: false,
  resolverData: {},
  resolverFetchError: undefined,
};

interface ResolverState {
  allSchemas: clutch.resolver.v1.Schema[];
  resolverData: object;
  resolverFetchError: ClutchError;
  resolverLoading: boolean;
  schemaFetchError: ClutchError;
  schemasLoading: boolean;
  searchableSchemas: clutch.resolver.v1.Schema[];
}

export interface DispatchAction {
  allSchemas?: any[];
  error?: ClutchError;
  schema?: any;
  type: ResolverAction;
}

const reducer = (state: ResolverState, action: DispatchAction) => {
  switch (action.type) {
    case ResolverAction.SCHEMAS_LOADING:
      return { ...initialState };
    case ResolverAction.SCHEMAS_SUCCCESS:
      return {
        ...state,
        schemasLoading: false,
        schemaFetchError: undefined,
        searchableSchemas: action.allSchemas
          .map(schema => {
            return schema.metadata.searchable || schema.metadata.search.enabled ? schema : null;
          })
          .filter(x => x),
        allSchemas: action.allSchemas,
      };
    case ResolverAction.SCHEMAS_ERROR:
      return {
        ...state,
        schemasLoading: false,
        schemaFetchError: action.error,
      };
    case ResolverAction.RESOLVING:
      return {
        ...state,
        resolverLoading: true,
        resolverFetchError: undefined,
      };
    case ResolverAction.RESOLVE_ERROR:
      return {
        ...state,
        resolverLoading: false,
        resolverFetchError: action.error,
      };
    case ResolverAction.RESOLVE_SUCCESS:
      return {
        ...state,
        resolverLoading: false,
        resolverFetchError: undefined,
      };
    default:
      throw new Error(`Unknown resolver action: ${action.type}`);
  }
};

const useResolverState = () => {
  return React.useReducer(reducer, initialState);
};

export { ResolverAction, useResolverState };
