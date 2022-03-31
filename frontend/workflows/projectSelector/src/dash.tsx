import * as React from "react";
import type { google as IGoogle } from "@clutch-sh/api";
import styled from "@emotion/styled";
import { Box } from "@material-ui/core";
import _ from "lodash";

import {
  ProjectSelectorDispatchContext,
  ProjectSelectorStateContext,
  TimelineDispatchContext,
  TimelineStateContext,
  TimeRangeDispatchContext,
  TimeRangeStateContext,
} from "./dash-hooks";
import { hoursToMs } from "./helpers";
import type { ProjectSelectorError } from "./project-selector";
import ProjectSelector from "./project-selector";
import type {
  DashAction,
  DashError,
  DashState,
  TimelineAction,
  TimelineState,
  TimeRangeAction,
  TimeRangeState,
} from "./types";

/**
 * DashProps: Defined input properties of the Dash component
 */
interface DashProps {
  /**
   * children: The children to render
   */
  children: React.ReactNode;
  /**
   * onError: (optional) error handler which will accept a DashError as input
   */
  onError?: (DashError) => void;
}

const initialState: DashState = {
  selected: [],
  projectData: {},
};

const initialTimelineState: TimelineState = {
  timeData: {},
};

const initialTimeRangeState: TimeRangeState = {
  // Default look back is 2 hours for time selector
  startTimeMs: Date.now() - hoursToMs(2),
  endTimeMs: Date.now(),
};

const CardContainer = styled.div({
  display: "flex",
  flex: 1,
  maxHeight: "100%",
  overflowY: "scroll",
});

const dashReducer = (state: DashState, action: DashAction): DashState => {
  switch (action.type) {
    case "UPDATE_SELECTED": {
      if (!_.isEqual(state.selected, action.payload.selected)) {
        return action.payload;
      }
      return state;
    }
    default:
      throw new Error("not implemented (should be unreachable)");
  }
};

const timelineReducer = (state: TimelineState, action: TimelineAction): TimelineState => {
  switch (action.type) {
    case "UPDATE": {
      // for now, clobber any existing data
      const newState = { ...state };
      newState.timeData[action.payload.key] = action.payload.eventData;
      return newState;
    }
    default:
      throw new Error("not implemented (should be unreachable)");
  }
};

const timeRangeReducer = (state: TimeRangeState, action: TimeRangeAction): TimeRangeState => {
  switch (action.type) {
    case "UPDATE": {
      // If start times do not match, or end times do not match, then return the update, else return the original state
      if (
        !_.isEqual(state.startTimeMs, action.payload.startTimeMs) ||
        !_.isEqual(state.endTimeMs, action.payload.endTimeMs)
      ) {
        return action.payload;
      }
      return state;
    }
    default:
      throw new Error("not implemented (should be unreachable)");
  }
};

const Dash = ({ children, onError }: DashProps) => {
  const [state, dispatch] = React.useReducer(dashReducer, initialState);
  const [timelineState, timelineDispatch] = React.useReducer(timelineReducer, initialTimelineState);
  const [timeRangeState, timeRangeDispatch] = React.useReducer(
    timeRangeReducer,
    initialTimeRangeState
  );

  // Will take a returned ProjectSelectorError and generate a rendered component of the error messages
  // and return it to the parent component
  const returnDashError = ({ errors }: ProjectSelectorError) => {
    if (onError && errors && errors.length) {
      const dashError: DashError = {
        title: "The following projects failed:",
      };
      dashError.data = (
        <ul>
          {errors.map((error: IGoogle.rpc.IStatus) => (
            <li>{error.message}</li>
          ))}
        </ul>
      );

      onError(dashError);
    }
  };

  return (
    <Box display="flex" flex={1} minHeight="100%" maxHeight="100%">
      {/* In an alternative design where the timeline drives the project selection, we could invert proj selector
       and timeline contexts. Right now the proj selector is the driver. */}
      <ProjectSelectorDispatchContext.Provider value={dispatch}>
        <ProjectSelectorStateContext.Provider value={state}>
          <TimeRangeDispatchContext.Provider value={timeRangeDispatch}>
            <TimeRangeStateContext.Provider value={timeRangeState}>
              <TimelineDispatchContext.Provider value={timelineDispatch}>
                <TimelineStateContext.Provider value={timelineState}>
                  <ProjectSelector onError={returnDashError} />
                  <CardContainer>{children}</CardContainer>
                </TimelineStateContext.Provider>
              </TimelineDispatchContext.Provider>
            </TimeRangeStateContext.Provider>
          </TimeRangeDispatchContext.Provider>
        </ProjectSelectorStateContext.Provider>
      </ProjectSelectorDispatchContext.Provider>
    </Box>
  );
};

export default Dash;
