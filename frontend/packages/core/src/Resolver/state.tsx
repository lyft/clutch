import React from "react";

import type { clutch } from "../../../../api";

enum ResolverAction {
  SCHEMAS_LOADING,
  SCHEMAS_SUCCCESS,
  SCHEMAS_ERROR,
  SET_SELECTED_SCHEMA,
  UPDATE_QUERY_DATA,
  RESOLVING,
  RESOLVE_ERROR,
  RESOLVE_SUCCESS,
}

const initialState = {
  schemasLoading: true,
  allSchemas: [],
  searchableSchemas: [],
  schemaFetchError: "",
  selectedSchema: 0,
  queryData: {},
  resolverLoading: false,
  resolverData: {},
  resolverFetchError: "",
};

interface ResolverState {
  allSchemas: clutch.resolver.v1.Schema[];
  queryData: {
    query: string;
  };
  resolverData: object;
  resolverFetchError: string;
  resolverLoading: boolean;
  schemaFetchError: string;
  schemasLoading: boolean;
  searchableSchemas: clutch.resolver.v1.Schema[];
  selectedSchema: number;
}

export interface DispatchAction {
  allSchemas?: any[];
  data?: {
    query?: string;
    [key: string]: string;
  };
  error?: string;
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
        schemaFetchError: "",
        searchableSchemas: action.allSchemas
          .map(schema => {
            return schema.metadata.searchable ? schema : null;
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
    case ResolverAction.SET_SELECTED_SCHEMA:
      return {
        ...state,
        selectedSchema: state.allSchemas.map(schema => schema.typeUrl).indexOf(action.schema),
        queryData: {},
        resolverFetchError: "",
      };
    case ResolverAction.UPDATE_QUERY_DATA: {
      const queryData = { ...state.queryData };
      if (action.data?.query === undefined) {
        delete queryData.query;
      }
      return {
        ...state,
        queryData: { ...queryData, ...action.data },
        resolverFetchError: "",
      };
    }
    case ResolverAction.RESOLVING:
      return {
        ...state,
        resolverLoading: true,
        resolverFetchError: "",
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
        resolverFetchError: "",
      };
    default:
      throw new Error(`Unknown resolver action: ${action.type}`);
  }
};

const useResolverState = () => {
  return React.useReducer(reducer, initialState);
};

export { ResolverAction, useResolverState };
