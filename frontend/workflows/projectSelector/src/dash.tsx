import * as React from "react";

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
    <DashDispatchContext.Provider value={dispatch}>
      <DashStateContext.Provider value={state}>
        <ProjectSelector />
        <div>{children}</div>
      </DashStateContext.Provider>
    </DashDispatchContext.Provider>
  );
};

export default Dash;
