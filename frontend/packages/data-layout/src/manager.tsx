import type { Thunk } from "react-hook-thunk-reducer";
import _ from "lodash";

import type { Action, ManagerLayout } from "./state";
import { ManagerAction, useManagerState } from "./state";

const assign = (key: string, value: object): Thunk<ManagerLayout, Action> => {
  return dispatch => {
    dispatch({
      type: ManagerAction.SET,
      payload: { key, value },
    });
  };
};

const update = (key: string, value: object): Thunk<ManagerLayout, Action> => {
  return dispatch => {
    dispatch({
      type: ManagerAction.UPDATE,
      payload: { key, value },
    });
  };
};

const hydrate = (key: string): Thunk<ManagerLayout, Action> => {
  return (dispatch, getState) => {
    const state = getState();
    if (Object.keys(state[key]).includes("hydrator")) {
      dispatch({
        type: ManagerAction.HYDRATE_START,
        payload: { key },
      });

      const args = state[key].deps.map(dep => state[dep].data);
      if (args.some(element => _.isEmpty(element))) {
        dispatch({
          type: ManagerAction.HYDRATE_END,
          payload: { key, error: `Missing depedency for data layout: ${key}` },
        });
        return;
      }

      return state[key]
        .hydrator(...args)
        .then(result => {
          dispatch({
            type: ManagerAction.HYDRATE_END,
            payload: {
              key,
              result: state[key].transformResponse(result),
            },
          });
        })
        .catch(error => {
          dispatch({
            type: ManagerAction.HYDRATE_END,
            payload: {
              key,
              error: state[key].transformError(error),
            },
          });
        });
    }
  };
};

interface Error {
  response?: {
    displayText?: string;
  };
  message: string;
}

interface DataManager {
  state: object;
  assign: (key: string, value: object) => void;
  hydrate: (key: string) => void;
  update: (key: string, value: object) => void;
}

const defaultTransform = (data: object): object => data;
const defaultErrorTransform = (err: Error): string => {
  return err?.response?.displayText ?? err.message;
};

const useDataLayoutManager = (layouts: ManagerLayout): DataManager => {
  const initialState = {};
  Object.keys(layouts).forEach(key => {
    const layout = layouts[key];
    initialState[key] = { data: {}, isLoading: true, error: null };
    if (layout?.hydrator !== undefined) {
      initialState[key] = {
        ...initialState[key],
        hydrator: layout?.hydrator || (() => {}),
        transformResponse: layout.transformResponse || defaultTransform,
        transformError: layout.transformError || defaultErrorTransform,
        deps: layout?.deps || [],
        cache: layout.cache ?? false,
      };
    }
  });

  const [state, dispatch] = useManagerState(initialState);
  return {
    state,
    assign: (key, value) => dispatch(assign(key, value)),
    hydrate: key => dispatch(hydrate(key)),
    update: (key, value) => dispatch(update(key, value)),
  };
};

export { DataManager, useDataLayoutManager };
