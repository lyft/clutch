import * as React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";

import type { DashAction, DashState, TimelineAction, TimelineState } from "./types";

export const DashStateContext = React.createContext<DashState | undefined>(undefined);
export const TimelineStateContext = React.createContext<TimelineState | undefined>(undefined);
export const TimelineUpdateContext = React.createContext<
  ((action: TimelineAction) => void) | undefined
>(undefined);
export const DashDispatchContext = React.createContext<((action: DashAction) => void) | undefined>(
  undefined
);

type useDashUpdaterReturn = {
  updateSelected: (state: DashState) => void;
};

type useTimelineUpdaterReturn = {
  updateTimeline: (key: string, points: IClutch.timeseries.v1.IPoint[]) => void;
};

export const useDashUpdater = (): useDashUpdaterReturn => {
  const dispatch = React.useContext(DashDispatchContext);

  return {
    updateSelected: projects => {
      dispatch && dispatch({ type: "UPDATE_SELECTED", payload: projects });
    },
  };
};

export const useTimelineUpdater = (): useTimelineUpdaterReturn => {
  const dispatch = React.useContext(TimelineUpdateContext);

  return {
    // TODO: how do we pass the key and points here instead of the state?
    updateTimeline: (key, points) => {
      dispatch && dispatch({ type: "UPDATE", payload: { key, points } });
    },
  };
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
