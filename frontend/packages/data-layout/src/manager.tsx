import type { Thunk } from "react-hook-thunk-reducer";
import type { ClutchError } from "@clutch-sh/core";
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

const reset = (): Thunk<ManagerLayout, Action> => {
  return dispatch => {
    dispatch({
      type: ManagerAction.RESET,
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
          payload: {
            key,
            error: {
              message: `Missing dependency for data layout: '${key}'`,
              status: {
                code: 404,
                text: "Not Found",
              },
            } as ClutchError,
          },
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

export interface DataManager {
  state: object;
  assign: (key: string, value: object) => void;
  hydrate: (key: string) => void;
  update: (key: string, value: object) => void;
  reset: () => void;
}

const defaultTransform = (data: object): object => data;
const defaultErrorTransform = (err: any): ClutchError => err;

const initialLayoutStepState = () => ({ data: {}, isLoading: true, error: null });

const defaultState = (layouts: ManagerLayout) => {
  const initializedLayouts = {};
  Object.keys(layouts).forEach(key => {
    const layout = layouts[key];
    initializedLayouts[key] = initialLayoutStepState();
    if (layout?.hydrator !== undefined) {
      initializedLayouts[key] = {
        ...initializedLayouts[key],
        hydrator: layout?.hydrator || (() => {}),
        transformResponse: layout.transformResponse || defaultTransform,
        transformError: layout.transformError || defaultErrorTransform,
        deps: layout?.deps || [],
        cache: layout.cache ?? false,
      };
    }
  });
  return initializedLayouts;
};

const useDataLayoutManager = (layouts: ManagerLayout): DataManager => {
  const initialState = defaultState(layouts);

  const [state, dispatch] = useManagerState(initialState);
  return {
    state,
    assign: (key, value) => dispatch(assign(key, value)),
    hydrate: key => dispatch(hydrate(key)),
    update: (key, value) => dispatch(update(key, value)),
    reset: () => dispatch(reset()),
  };
};

export default useDataLayoutManager;
