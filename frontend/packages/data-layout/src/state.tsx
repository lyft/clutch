import type React from "react";
import type { Thunk } from "react-hook-thunk-reducer";
import useThunkReducer from "react-hook-thunk-reducer";
import type { ClutchError } from "@clutch-sh/core";

enum ManagerAction {
  HYDRATE_START,
  HYDRATE_END,
  SET,
  UPDATE,
  RESET,
}

export interface ManagerLayout {
  [key: string]: {
    isLoading?: boolean;
    data?: object;
    error?: ClutchError;
    hydrator?: (...args: any[]) => any;
    transformResponse?: (...args: any[]) => any;
    transformError?: (...args: any[]) => ClutchError;
    deps?: string[];
    cache?: boolean;
  };
}

interface ActionPayload {
  key?: string;
  value?: object;
  result?: object;
  error?: ClutchError;
  overwrite?: boolean;
}

export interface Action {
  type: ManagerAction;
  payload?: ActionPayload;
}

const reducer = (state: ManagerLayout, action: Action): ManagerLayout => {
  const layoutKey = action?.payload?.key;
  const stateClone = { ...state };

  switch (action.type) {
    case ManagerAction.HYDRATE_START:
      return {
        ...state,
        [layoutKey]: { ...state[layoutKey], isLoading: true },
      };
    case ManagerAction.HYDRATE_END: {
      const newData: any = action.payload?.result;
      const existingData: any = state[layoutKey]?.data;
      const newDataIsArray = Array.isArray(newData);
      const existingDataIsArray = Array.isArray(existingData);
      let data: object;
      if (((newDataIsArray && !existingDataIsArray) || (!newDataIsArray && existingDataIsArray)) || action.payload?.overwrite) {
        data = newData;
      } else if (newDataIsArray) {
        data = [...existingData, ...newData];
      } else {
        data = { ...existingData, ...newData };
      }
      const update = {
        isLoading: false,
        data,
        error: action.payload?.error,
      };
      return {
        ...state,
        [layoutKey]: { ...state[layoutKey], ...update },
      };
    }
    case ManagerAction.SET:
      return {
        ...state,
        [layoutKey]: {
          ...state?.[layoutKey],
          data: action.payload?.value,
          isLoading: false,
        },
      };
    case ManagerAction.UPDATE:
      return {
        ...state,
        [layoutKey]: {
          isLoading: false,
          ...state[layoutKey],
          ...action.payload?.value,
        },
      };
    case ManagerAction.RESET:
      Object.keys(state).forEach(key => {
        stateClone[key] = {
          ...state[key],
          data: {},
          isLoading: true,
          error: null,
        };
      });
      return stateClone;
    default:
      throw new Error(`Unknown data manager action: ${action.type}`);
  }
};

const useManagerState = (
  initialState: ManagerLayout
): [ManagerLayout, React.Dispatch<Action | Thunk<ManagerLayout, Action>>] => {
  return useThunkReducer(reducer, initialState);
};

export { ManagerAction, reducer, useManagerState };
