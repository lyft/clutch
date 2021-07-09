import * as React from "react";
import ProjectSelector from "./project-selector";
import type { Group } from "./types";

const useCardManager = () => {};

// update selected projects
export const useDashManager = () => {};

type DashActionKind = "UPDATE_PROJECTS";

interface DashAction {
  type: DashActionKind;
  payload: string[];
}

interface DashState {
  selectedProjects: string[];
}

const StateContext = React.createContext<DashState | undefined>(undefined);
const useReducerState = () => {
  return React.useContext(StateContext);
};

const DispatchContext = React.createContext<(action: DashAction) => void | undefined>(undefined);
const useDispatch = () => {
  return React.useContext(DispatchContext);
};

type useProjectUpdaterReturn = {
  updateProjects: (projects: string[]) => void;
};

export const useProjectUpdater = (): useProjectUpdaterReturn => {
  const dispatch = React.useContext(DispatchContext);

  return {
    updateProjects: projects => {
      dispatch({ type: "UPDATE_PROJECTS", payload: projects });
    },
  };
};

const initialState = {
  selectedProjects: [],
};

const dashReducer = (state: DashState, action: DashAction): DashState => {
  switch (action.type) {
    case "UPDATE_PROJECTS": {
      return {...state, selectedProjects: action.payload };
    }
    default:
      throw new Error("not implemented (should be unreachable)");
  }
};

const Card = () => {
  const {selectedProjects} = useReducerState();

  return <div>Hello world! {JSON.stringify(selectedProjects)}</div>;
};

export const Dash = () => {
  const [state, dispatch] = React.useReducer(dashReducer, initialState);

  return (
    <DispatchContext.Provider value={dispatch}>
      <StateContext.Provider value={state}>
        <ProjectSelector />
        <div>
          <Card />
        </div>
      </StateContext.Provider>
    </DispatchContext.Provider>
  );
};
