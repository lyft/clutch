import * as React from "react";
import ProjectSelector from "./project-selector";
import type { DashState } from "./types";

type DashActionKind = "UPDATE_SELECTED";

interface DashAction {
  type: DashActionKind;
  payload: DashState;
}

const StateContext = React.createContext<DashState | undefined>(undefined);

export const useDashState = () => {
  return React.useContext(StateContext);
};

const DispatchContext = React.createContext<(action: DashAction) => void | undefined>(undefined);

type useDashUpdater = {
  updateSelected: (state: DashState) => void;
};

export const useDashUpdater = (): useDashUpdater => {
  const dispatch = React.useContext(DispatchContext);

  return {
    updateSelected: projects => {
      dispatch({ type: "UPDATE_SELECTED", payload: projects });
    },
  };
};

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

export const Dash = ({children}) => {
  const [state, dispatch] = React.useReducer(dashReducer, initialState);

  return (
    <DispatchContext.Provider value={dispatch}>
      <StateContext.Provider value={state}>
        <ProjectSelector />
        <div>
          {children}
        </div>
      </StateContext.Provider>
    </DispatchContext.Provider>
  );
};
