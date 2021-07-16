import * as React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import type { ClutchError } from "@clutch-sh/core";
import { client, TextField, userId } from "@clutch-sh/core";
import styled from "@emotion/styled";
import { Divider, LinearProgress } from "@material-ui/core";
import LayersIcon from "@material-ui/icons/Layers";
import _ from "lodash";

import { useDashUpdater } from "./dash-hooks";
import { DispatchContext, StateContext } from "./helpers";
import ProjectGroup from "./project-group";
import selectorReducer from "./selector-reducer";
import type { DashState, State } from "./types";
import { Group } from "./types";

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

  const { updateSelected } = useDashUpdater();

  const [state, dispatch] = React.useReducer(selectorReducer, initialState);

  React.useEffect(() => {
    console.log("effect"); // eslint-disable-line
    // Determine if any hydration is required.
    // - Are any services missing from state.projectdata?
    // - Are projects empty (first load)?
    // - Is loading not already in progress?

    let allPresent = true;
    _.forEach(Object.keys(state[Group.PROJECTS]), p => {
      /*
      TODO: b/c of this conditional, if a user adds an upstream/downstream we already have the project data for
      to the custom project group, allPresent will be true and we wont trigger an api call. One way to account for this
      is updating the conditional to additionally check if the project is included in state[Group.Downstreams]/state[Group.Upstreams]
      and if so, mark allPresent as false.
      */
      if (!(p in state.projectData)) {
        allPresent = false;
        return false; // Stop iteration.
      }
      return true; // Continue.
    });

    if (!state.loading && (Object.keys(state[Group.PROJECTS]).length === 0 || !allPresent)) {
      console.log("calling API!", state.loading); // eslint-disable-line
      dispatch({ type: "HYDRATE_START" });

      // TODO: have userId check be server driven
      const requestParams = { users: [userId()], projects: [] };
      _.forEach(Object.keys(state[Group.PROJECTS]), p => {
        // if the project is custom
        if (state[Group.PROJECTS][p].custom) {
          requestParams.projects.push(p);
        }
      });

      client
        .post("/v1/project/getProjects", requestParams as IClutch.project.v1.GetProjectsRequest)
        .then(resp => {
          const { results } = resp.data as IClutch.project.v1.GetProjectsResponse;
          dispatch({ type: "HYDRATE_END", payload: { result: results || {} } });
        })
        .catch((err: ClutchError) => {
          dispatch({ type: "HYDRATE_ERROR", payload: { result: err } });
        });
    }
  }, [state[Group.PROJECTS]]);

  // This hook updates the global dash state based on the currently selected projects for cards to consume (including upstreams and downstreams).
  React.useEffect(() => {
    const dashState: DashState = { projectData: {}, selected: [] };

    // Determine selected projects.
    const selected = new Set<string>();
    _.forEach(Object.keys(state[Group.PROJECTS]), p => {
      if (state[Group.PROJECTS][p].checked) {
        selected.add(p);
      }
    });
    _.forEach(Object.keys(state[Group.DOWNSTREAM]), p => {
      if (state[Group.DOWNSTREAM][p].checked) {
        selected.add(p);
      }
    });
    _.forEach(Object.keys(state[Group.UPSTREAM]), p => {
      if (state[Group.UPSTREAM][p].checked) {
        selected.add(p);
      }
    });
    dashState.selected = Array.from(selected).sort();

    // Collect project data.
    _.forEach(dashState.selected, p => {
      dashState.projectData[p] = state.projectData[p];
    });

    // Update!
    updateSelected(dashState);
  }, [state]);

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
            helperText={state.error?.message}
            error={hasError}
          />
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

export default ProjectSelector;
