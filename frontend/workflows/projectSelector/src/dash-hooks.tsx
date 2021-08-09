import * as React from "react";

import type { DashAction, DashState } from "./types";

export const DashStateContext = React.createContext<DashState | undefined>(undefined);

export const DashDispatchContext = React.createContext<((action: DashAction) => void) | undefined>(
  undefined
);

type useDashUpdaterReturn = {
  updateSelected: (state: DashState) => void;
};

export const useDashUpdater = (): useDashUpdaterReturn => {
  const dispatch = React.useContext(DashDispatchContext);

  return {
    updateSelected: projects => {
      dispatch && dispatch({ type: "UPDATE_SELECTED", payload: projects });
    },
  };
};

export const useDashState = (): DashState => {
  const value = React.useContext<DashState | undefined>(DashStateContext);
  if (!value) {
    throw new Error(
      "useDashState was invoked with no value, check that it is a child of the Dash component"
    );
  }
  return value;
};
