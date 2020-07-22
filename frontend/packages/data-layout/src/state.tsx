import type React from "react";
import type { Thunk } from "react-hook-thunk-reducer";
import useThunkReducer from "react-hook-thunk-reducer";

enum ManagerAction {
  HYDRATE_START,
  HYDRATE_END,
  SET,
  UPDATE,
}

export interface ManagerLayout {
  [key: string]: {
    isLoading?: boolean;
    data?: object;
    error?: string;
    hydrator?: (...args: any[]) => any;
    transformResponse?: (...args: any[]) => any;
    transformError?: (...args: any[]) => any;
    deps?: string[];
    cache?: boolean;
  };
}

interface ActionPayload {
  key?: string;
  value?: object;
  result?: object;
  error?: string;
}

export interface Action {
  type: ManagerAction;
  payload?: ActionPayload;
}

const reducer = (state: ManagerLayout, action: Action): ManagerLayout => {
  const layoutKey = action?.payload?.key;

  switch (action.type) {
    case ManagerAction.HYDRATE_START:
      return {
        ...state,
        [layoutKey]: { ...state[layoutKey], isLoading: true },
      };
    case ManagerAction.HYDRATE_END: {
      const update = {
        isLoading: false,
        data: action.payload?.result || {},
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
          ...state[layoutKey],
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
    default:
      throw new Error(`Unknown data manager action: ${action.type}`);
  }
};

const useManagerState = (
  initialState: ManagerLayout
): [ManagerLayout, React.Dispatch<Action | Thunk<ManagerLayout, Action>>] => {
  return useThunkReducer(reducer, initialState);
};

export { ManagerAction, useManagerState };
