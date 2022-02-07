import * as React from "react";
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
import ProjectSelector from "./project-selector";
import type { ProjectSelectorError } from "./project-selector";
import type {
  DashAction,
  DashState,
  TimelineAction,
  TimelineState,
  TimeRangeAction,
  TimeRangeState,
} from "./types";

export interface DashError {
  title: string;
  message: string;
  data?: React.ReactNode;
}

interface DashProps {
  children: React.ReactNode;
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
    // TODO: Add more actions like slicing by time
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
      // TODO: when adding more complexity to the time range state - i.e. filters of event types - add logic here to handle if both start and end times are equal
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

  const returnDashError = ({ projects, type }: ProjectSelectorError) => {
    if (onError && type === "DEPRECATED") {
      const errorData: DashError = {
        title: "Project Deprecation",
        message:
          "The following projects have been deprecated and removed from the current selections:",
      };

      if (projects && projects.length) {
        errorData.data = (
          <ul>
            {projects.map(project => (
              <li>{project}</li>
            ))}
          </ul>
        );
      }

      onError(errorData);
    }
  };

  return (
    <Box display="flex" flex={1} minHeight="100%" maxHeight="100%">
      {/* TODO: Maybe in the future invert proj selector and timeline contexts */}
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
