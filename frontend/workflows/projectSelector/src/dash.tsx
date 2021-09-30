import * as React from "react";
import styled from "@emotion/styled";
import { Box } from "@material-ui/core";
import _ from "lodash";

import {
  ProjectSelectorDispatchContext,
  ProjectSelectorStateContext,
  TimelineDispatchContext,
  TimelineStateContext,
} from "./dash-hooks";
import ProjectSelector from "./project-selector";
import type { DashAction, DashState, TimelineAction, TimelineState } from "./types";

const initialState = {
  selected: [],
  projectData: {},
};

const initialTimelineState = {
  timeData: {},
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
      newState.timeData[action.payload.key] = action.payload.points;
      return newState;
    }
    default:
      throw new Error("not implemented (should be unreachable)");
  }
};

const Dash = ({ children }) => {
  const [state, dispatch] = React.useReducer(dashReducer, initialState);
  const [timelineState, timelineDispatch] = React.useReducer(timelineReducer, initialTimelineState);
  return (
    <Box display="flex" flex={1} minHeight="100%" maxHeight="100%">
      {/* TODO: Maybe in the future invert proj selector and timeline contexts */}
      <ProjectSelectorDispatchContext.Provider value={dispatch}>
        <ProjectSelectorStateContext.Provider value={state}>
          <TimelineDispatchContext.Provider value={timelineDispatch}>
            <TimelineStateContext.Provider value={timelineState}>
              <ProjectSelector />
              <CardContainer>{children}</CardContainer>
            </TimelineStateContext.Provider>
          </TimelineDispatchContext.Provider>
        </ProjectSelectorStateContext.Provider>
      </ProjectSelectorDispatchContext.Provider>
    </Box>
  );
};

export default Dash;
