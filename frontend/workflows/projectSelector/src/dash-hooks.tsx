import * as React from "react";

import type { DashAction, DashState, TimelineState } from "./types";

export const DashStateContext = React.createContext<DashState | undefined>(undefined);
export const TimelineStateContext = React.createContext<TimelineState | undefined>(undefined);
// TODO: What type goes here? React.Dispatch<any> or? ???
export const TimelineUpdateContext = React.createContext<React.Dispatch<any> | undefined>(
  undefined
);
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

export const useTimelineUpdate = (): React.Dispatch<any> => {
  const setTimeData = React.useContext(TimelineUpdateContext);
  if (!setTimeData) {
    throw new Error(
      "useTimelineUpdate was invoked outside of a valid context, check that it is a child of the Timeline component"
    );
  }
  return setTimeData;
};

export const useDashState = (): DashState => {
  const value = React.useContext<DashState | undefined>(DashStateContext);
  if (!value) {
    throw new Error(
      "useDashState was invoked outside of a valid context, check that it is a child of the Dash component"
    );
  }
  return value;
};

export const useTimelineState = (): TimelineState => {
  const value = React.useContext<TimelineState | undefined>(TimelineStateContext);
  if (!value) {
    throw new Error(
      "useTimelineState was invoked outside of a valid context, check that it is a child of the Timeline component"
    );
  }
  return value;
};
