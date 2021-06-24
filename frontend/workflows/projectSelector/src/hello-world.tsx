import * as React from "react";
import type { ClutchError } from "@clutch-sh/core";
import { client, Error, TextField, userId } from "@clutch-sh/core";
import styled from "@emotion/styled";
import { Divider, LinearProgress } from "@material-ui/core";
import LayersIcon from "@material-ui/icons/Layers";
import _ from "lodash";

import ProjectGroup from "./project-group";
import { selectorReducer } from "./selector-reducer";

export enum Group {
  PROJECTS,
  UPSTREAM,
  DOWNSTREAM,
}

type UserActionKind =
  | "ADD_PROJECTS"
  | "REMOVE_PROJECTS"
  | "TOGGLE_PROJECTS"
  | "TOGGLE_ENTIRE_GROUP"
  | "ONLY_PROJECTS";

interface UserAction {
  type: UserActionKind;
  payload: UserPayload;
}

interface UserPayload {
  group: Group;
  projects?: string[];
}

type BackgroundActionKind = "HYDRATE_START" | "HYDRATE_END" | "HYDRATE_ERROR";

interface BackgroundAction {
  type: BackgroundActionKind;
  payload?: BackgroundPayload;
}

interface BackgroundPayload {
  result: any;
}

export type Action = BackgroundAction | UserAction;

export interface State {
  [Group.PROJECTS]: GroupState;
  [Group.UPSTREAM]: GroupState;
  [Group.DOWNSTREAM]: GroupState;

  projectData: { [projectName: string]: Project };
  loading: boolean;
  error: ClutchError | undefined;
}

interface Project {
  name: string;
  tier: string;
  owners: string[];
  languages: string[];
  data: any;
  upstreams: string[];
  downstreams: string[];
}

interface GroupState {
  [projectName: string]: ProjectState;
}

interface ProjectState {
  checked: boolean;
  // TODO: hidden should be derived?
  hidden?: boolean; // upstreams and downstreams are hidden when their parent is unchecked unless other parents also use them.
  custom?: boolean;
}

const StateContext = React.createContext<State | undefined>(undefined);
export const useReducerState = () => {
  return React.useContext(StateContext);
};

const DispatchContext = React.createContext<(action: Action) => void | undefined>(undefined);
export const useDispatch = () => {
  return React.useContext(DispatchContext);
};

// TODO(perf): call with useMemo().
export const deriveSwitchStatus = (state: State, group: Group): boolean => {
  return (
    Object.keys(state[group]).length > 0 &&
    Object.keys(state[group]).every(key => state[group][key].checked)
  );
};

const initialState: State = {
  [Group.PROJECTS]: {},
  [Group.UPSTREAM]: {},
  [Group.DOWNSTREAM]: {},
  projectData: {},
  loading: false,
  error: undefined,
};

const StyledSelectorContainer = styled.div({
  backgroundColor: "#F9FAFE",
  borderRight: "1px solid rgba(13, 16, 48, 0.1)",
  boxShadow: "0px 4px 6px rgba(53, 72, 212, 0.2)",
  width: "245px",
});

const StyledWorkflowHeader = styled.div({
  margin: "16px 16px 12px 16px",
  display: "flex",
  alignItems: "center",
});

const StyledWorkflowTitle = styled.span({
  fontWeight: "bold",
  fontSize: "20px",
  lineHeight: "24px",
  margin: "0px 8px",
});

const StyledProjectTextField = styled(TextField)({
  padding: "16px 16px 8px 16px",
});

const StyledProgressContainer = styled.div({
  height: "4px",
  ".MuiLinearProgress-root": {
    backgroundColor: "rgb(194, 200, 242)",
  },
  ".MuiLinearProgress-bar": {
    backgroundColor: "#3548D4",
  },
});

const ProjectSelector = () => {
  // On load, we'll request a list of owned projects and their upstreams and downstreams from the API.
  // The API will contain information about the relationships between projects and upstreams and downstreams.
  // By default, the owned projects will be checked and others will be unchecked.
  // TODO: If a project is unchecked, the upstream and downstreams related to it disappear from the list.
  // TODO: If a project is rechecked, the checks were preserved.

  const [customProject, setCustomProject] = React.useState("");

  const [state, dispatch] = React.useReducer(selectorReducer, initialState);

  React.useEffect(() => {
    console.log("effect");
    // Determine if any hydration is required.
    // - Are any services missing from state.projectdata?
    // - Are projects empty (first load)?
    // - Is loading not already in progress?

    let allPresent = true;
    _.forEach(Object.keys(state[Group.PROJECTS]), p => {
      /*
      TODO: b/c of this conditional, if a user adds an upstream/downstream we already have the project data for
      to the custom project group, allPresent will be true and we wont trigger an api call
      */
      if (!(p in state.projectData)) {
        allPresent = false;
        return false; // Stop iteration.
      }
      return true; // Continue.
    });

    if (!state.loading && (Object.keys(state[Group.PROJECTS]).length == 0 || !allPresent)) {
      console.log("calling API!", state.loading);
      dispatch({ type: "HYDRATE_START" });

      // TODO: have userId check be server driven
      const requestParams = {"users": [userId()], "projects": []}

      const customProjects = [];
      _.forEach(Object.keys(state[Group.PROJECTS]), p => {
        // if the project is custom and missing from state.projectdata
        if (state[Group.PROJECTS][p].custom && !(p in state.projectData)) {
          customProjects.push(p);
        }
      });
      if (customProjects.length > 0) {
        requestParams.projects = customProjects;
      }

      /*
      TODO: the API doesn't return an error if a custom project is not found so we should first
      check if the API returns empty results and process that as an error
      */
      client
        .post("/v1/project/getProjects", requestParams)
        .then(resp => dispatch({ type: "HYDRATE_END", payload: { result: resp.data.results } }))
        .catch((err: ClutchError) => {
          dispatch({ type: "HYDRATE_ERROR", payload: { result: err } });
        });

    }
  }, [state[Group.PROJECTS]]);

  const handleAdd = () => {
    if (customProject === "") {
      return;
    }
    dispatch({
      type: "ADD_PROJECTS",
      payload: { group: Group.PROJECTS, projects: [customProject] },
    });
    setCustomProject("");
  };

  const hasError = state.error !== undefined && state.error !== null;

  return (
    <DispatchContext.Provider value={dispatch}>
      <StateContext.Provider value={state}>
        <StyledSelectorContainer>
          <StyledWorkflowHeader>
            {/* TODO: change icon to match design */}
            <LayersIcon />
            <StyledWorkflowTitle>Dash</StyledWorkflowTitle>
          </StyledWorkflowHeader>
          <StyledProgressContainer>
            {state.loading && <LinearProgress color="secondary" />}
          </StyledProgressContainer>
          <Divider />
          {/* TODO: add plus icon in the text field */}
          <StyledProjectTextField
            disabled={state.loading}
            placeholder="Add a project"
            value={customProject}
            onChange={e => setCustomProject(e.target.value)}
            onKeyDown={e => e.key === "Enter" && handleAdd()}
          />
          {/* TODO: styling for the error */}
          {hasError && <Error subject={state.error} />}
          <ProjectGroup title="Projects" group={Group.PROJECTS} displayToggleHelperText />
          <Divider />
          <ProjectGroup title="Upstreams" group={Group.UPSTREAM} />
          <Divider />
          <ProjectGroup title="Downstreams" group={Group.DOWNSTREAM} />
        </StyledSelectorContainer>
      </StateContext.Provider>
    </DispatchContext.Provider>
  );
};

const HelloWorld = () => (
  <>
    <ProjectSelector />
  </>
);

export default HelloWorld;
