import * as React from "react";

import type { DashAction, DashState, TimeDataUpdate, TimelineAction, TimelineState } from "./types";

// Contexts for project selector
export const ProjectSelectorStateContext = React.createContext<DashState | undefined>(undefined);
export const ProjectSelectorDispatchContext = React.createContext<
  ((action: DashAction) => void) | undefined
>(undefined);

// Contexts for timeline
export const TimelineStateContext = React.createContext<TimelineState | undefined>(undefined);
export const TimelineDispatchContext = React.createContext<
  ((action: TimelineAction) => void) | undefined
>(undefined);

// project selector hooks
type useDashUpdaterReturn = {
  updateSelected: (state: DashState) => void;
};

export const useDashUpdater = (): useDashUpdaterReturn => {
  const dispatch = React.useContext(ProjectSelectorDispatchContext);

  return {
    updateSelected: projects => {
      dispatch && dispatch({ type: "UPDATE_SELECTED", payload: projects });
    },
  };
};

export const useDashState = (): DashState => {
  const value = React.useContext<DashState | undefined>(ProjectSelectorStateContext);
  if (!value) {
    throw new Error(
      "useDashState was invoked outside of a valid context, check that it is a child of the Dash component"
    );
  }
  return value;
};

// timeline hooks
type useTimelineUpdaterReturn = {
  updateTimeline: (update: TimeDataUpdate) => void;
};

// hook for writing
export const useTimelineUpdater = (): useTimelineUpdaterReturn => {
  const dispatch = React.useContext(TimelineDispatchContext);

  return {
    updateTimeline: (update: TimeDataUpdate) => {
      dispatch && dispatch({ type: "UPDATE", payload: update });
    },
  };
};

// hook for reading
export const useTimelineState = (): TimelineState => {
  const value = React.useContext<TimelineState | undefined>(TimelineStateContext);
  if (!value) {
    throw new Error(
      "useTimelineState was invoked outside of a valid context, check that it is a child of the Timeline component"
    );
  }
  return value;
};
