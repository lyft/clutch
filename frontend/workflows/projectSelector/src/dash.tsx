import * as React from "react";
import { Box } from "@material-ui/core";
import _ from "lodash";

import { DashDispatchContext, DashStateContext } from "./dash-hooks";
import ProjectSelector from "./project-selector";
import type { DashAction, DashState } from "./types";

const initialState = {
  selected: [],
  projectData: {},
  delta: {
    newSelected: [],
    newUnselected: [],
  },
};

const dashReducer = (state: DashState, action: DashAction): DashState => {
  switch (action.type) {
    case "UPDATE_SELECTED": {
      // comparisons of what was selected and removed
      const selected = _.difference(action.payload.selected, state.selected);
      const unselected = _.difference(state.selected, action.payload.selected);

      if (!_.isEqual(state.selected, action.payload.selected)) {
        return {
          selected: action.payload.selected,
          projectData: action.payload.projectData,
          delta: {
            newSelected: selected,
            newUnselected: unselected,
          },
        };
      }
      return state;
    }
    default:
      throw new Error("not implemented (should be unreachable)");
  }
};

const Dash = ({ children }) => {
  const [state, dispatch] = React.useReducer(dashReducer, initialState);

  return (
    <Box display="flex" flex={1}>
      <DashDispatchContext.Provider value={dispatch}>
        <DashStateContext.Provider value={state}>
          <ProjectSelector />
          <Box display="flex" flex={1}>
            {children}
          </Box>
        </DashStateContext.Provider>
      </DashDispatchContext.Provider>
    </Box>
  );
};

export default Dash;
