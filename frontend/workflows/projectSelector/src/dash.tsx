import * as React from "react";
import ProjectSelector from "./project-selector";

// update selected projects
export const useDashManager = () => {};

type DashActionKind = "UPDATE_PROJECTS";

interface DashAction {
  type: DashActionKind;
  payload: string[];
}

interface DashState {
  // TODO: richer state, including upstreams, downstreams, and full project data.
  selectedProjects: string[];

  // TODO: add events to dash state or publish to separate events context? consider pros/cons (possible rendering tax).
}

const StateContext = React.createContext<DashState | undefined>(undefined);

// TODO: split out into smaller hooks that don't return full state.
const useReducerState = () => {
  return React.useContext(StateContext);
};

const DispatchContext = React.createContext<(action: DashAction) => void | undefined>(undefined);

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
      return { ...state, selectedProjects: action.payload };
    }
    default:
      throw new Error("not implemented (should be unreachable)");
  }
};

const Card = () => {
  const { selectedProjects } = useReducerState();

  return (
    <div style={{ padding: "10px", margin: "10px", border: "1px solid grey", borderRadius: "4px" }}>
      Hello world card!
      <br /> Projects: {JSON.stringify(selectedProjects)}
    </div>
  );
};

// TODO: ability to add cards via children.
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
