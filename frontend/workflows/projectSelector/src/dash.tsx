import * as React from "react";
import styled from "@emotion/styled";
import { Box, Fab, useMediaQuery } from "@material-ui/core";
import KeyboardArrowDownIcon from "@material-ui/icons/KeyboardArrowDown";
import KeyboardArrowUpIcon from "@material-ui/icons/KeyboardArrowUp";
import _ from "lodash";

import { DashDispatchContext, DashStateContext } from "./dash-hooks";
import ProjectSelector from "./project-selector";
import type { DashAction, DashState } from "./types";

const initialState = {
  selected: [],
  projectData: {},
};

const StyledFab = styled(Fab)({
  position: "fixed",
  left: "85px",
  bottom: "50vh",
  transform: "rotate(270deg)",
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

const Dash = ({ children }) => {
  const [state, dispatch] = React.useReducer(dashReducer, initialState);
  const compress = useMediaQuery((theme: any) => theme.breakpoints.down("md"));
  const [showSelector, setShowSelector] = React.useState<boolean>(false);
  const Icon = showSelector ? KeyboardArrowUpIcon : KeyboardArrowDownIcon;

  return (
    <Box display="flex" flex={1}>
      <DashDispatchContext.Provider value={dispatch}>
        <DashStateContext.Provider value={state}>
          {compress ? (
            <>
              <StyledFab onClick={() => setShowSelector(s => !s)} size="small">
                <Icon />
              </StyledFab>
              {showSelector ? (
                <ProjectSelector />
              ) : (
                <Box display="flex" flex={1}>
                  {children}
                </Box>
              )}
            </>
          ) : (
            <>
              <ProjectSelector />
              <Box display="flex" flex={1}>
                {children}
              </Box>
            </>
          )}
        </DashStateContext.Provider>
      </DashDispatchContext.Provider>
    </Box>
  );
};

export default Dash;
