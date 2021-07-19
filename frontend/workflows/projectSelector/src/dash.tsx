import * as React from "react";
import { Box } from "@material-ui/core";

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
      return action.payload;
    }
    default:
      throw new Error("not implemented (should be unreachable)");
  }
};

const Dash = ({ children }) => {
  const [state, dispatch] = React.useReducer(dashReducer, initialState);

  return (
    <Box display="flex" width="100%">
      <DashDispatchContext.Provider value={dispatch}>
        <DashStateContext.Provider value={state}>
          <ProjectSelector />
          <Box flexGrow={1}>{children}</Box>
        </DashStateContext.Provider>
      </DashDispatchContext.Provider>
    </Box>
  );
};

export default Dash;
