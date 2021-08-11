import * as React from "react";
import styled from "@emotion/styled";
import { Box, Fab, useMediaQuery } from "@material-ui/core";
import ChevronRightIcon from "@material-ui/icons/ChevronRight";
import _ from "lodash";

import { DashDispatchContext, DashStateContext } from "./dash-hooks";
import ProjectSelector from "./project-selector";
import type { DashAction, DashState } from "./types";

const initialState = {
  selected: [],
  projectData: {},
};

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

const StyledFab = styled(Fab)({
  position: "fixed",
  left: "80px",
  top: "70px",
});

const Dash = ({ children }) => {
  const [state, dispatch] = React.useReducer(dashReducer, initialState);
  const [selectorOpen, setSelectorOpen] = React.useState<boolean>(false);
  const compressed = useMediaQuery((theme: any) => theme.breakpoints.down("md"));

  return (
    <Box display="flex" flex={1}>
      <DashDispatchContext.Provider value={dispatch}>
        <DashStateContext.Provider value={state}>
          {compressed && !selectorOpen && (
            <StyledFab size="small" onClick={() => setSelectorOpen(true)}>
              <ChevronRightIcon style={{ marginLeft: "10px" }} />
            </StyledFab>
          )}
          {(!compressed || (compressed && selectorOpen)) && (
            <ProjectSelector fullScreen={compressed} onClose={() => setSelectorOpen(false)} />
          )}
          {(!compressed || (compressed && !selectorOpen)) && (
            <Box display="flex" flex={1}>
              {children}
            </Box>
          )}
        </DashStateContext.Provider>
      </DashDispatchContext.Provider>
    </Box>
  );
};

export default Dash;
