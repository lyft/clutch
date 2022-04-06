import React from "react";
import { client, Grid, Paper, TextField, Typography, useNavigate } from "@clutch-sh/core";
import { Box, CircularProgress } from "@material-ui/core";
import SearchIcon from "@material-ui/icons/Search";

import type { WorkflowProps } from "..";

import catalogReducer from "./catalog-reducer";
import ProjectCard from "./project-card";
import { addProject, getProjects, removeProject } from "./storage";
import type { CatalogState } from "./types";

const initialState: CatalogState = {
  projects: [],
  search: "",
  isLoading: false,
  isSearching: false,
  error: undefined,
};

const Catalog: React.FC<WorkflowProps> = ({ heading }) => {
  const navigate = useNavigate();
  const [state, dispatch] = React.useReducer(catalogReducer, initialState);

  const navigateToProject = project => {
    navigate(`/catalog/${project.name}`);
  };

  const setError = err => dispatch({ type: "HYDRATE_ERROR", payload: { result: err.message } });

  React.useEffect(() => {
    dispatch({ type: "HYDRATE_START" });
    getProjects(
      projects => dispatch({ type: "HYDRATE_END", payload: { result: projects } }),
      setError
    );
  }, []);

  // TODO: Decouple some of the logic in the storage functions and migrate it to the reducer
  const triggerProjectAdd = () => {
    dispatch({ type: "SEARCH_START" });
    addProject(
      state?.search || "",
      projects => {
        dispatch({ type: "ADD_PROJECT", payload: { projects } });
        dispatch({ type: "SEARCH_END" });
      },
      e => {
        dispatch({ type: "SEARCH_END" });
        setError(e);
      }
    );
    const projectMatches = state.projects.filter(
      p => state?.search && state.search !== "" && p?.name === state.search
    );
    if (projectMatches.length === 1) {
      navigateToProject(projectMatches[0]);
    }
  };

  const triggerProjectRemove = project => {
    removeProject(
      project.name,
      projects => dispatch({ type: "REMOVE_PROJECT", payload: { projects } }),
      setError
    );
  };

  const autoComplete = async (search: string): Promise<any> => {
    // Check the length of the search query as the user might empty out the search
    // which will still trigger the on change handler
    if (search.length === 0) {
      return { results: [] };
    }

    const response = await client.post("/v1/resolver/autocomplete", {
      want: `type.googleapis.com/clutch.core.project.v1.Project`,
      search,
    });

    return { results: response?.data?.results || [] };
  };

  return (
    <Box style={{ padding: "32px" }}>
      <div style={{ marginBottom: "8px" }}>
        <Typography variant="caption2" color="rgb(13, 16, 48, .48)">
          Project Catalog&nbsp;/&nbsp;Index
        </Typography>
      </div>
      <div style={{ marginBottom: "32px" }}>
        <Typography variant="h2">Project Catalog</Typography>
      </div>
      <Paper>
        <div style={{ margin: "16px" }}>
          <TextField
            placeholder="Search"
            value={state.search}
            onChange={e => dispatch({ type: "SEARCH", payload: { search: e.target.value } })}
            onKeyDown={e => e.key === "Enter" && triggerProjectAdd()}
            autocompleteCallback={v => autoComplete(v)}
            endAdornment={
              state.isSearching ? (
                <CircularProgress size="24px" />
              ) : (
                <SearchIcon onClick={triggerProjectAdd} />
              )
            }
            error={state.error !== undefined}
            helperText={state?.error}
          />
        </div>
      </Paper>
      <div style={{ marginBottom: "16px", marginTop: "32px" }}>
        <Typography variant="h3">My Projects</Typography>
      </div>
      <Grid container direction="row" spacing={5}>
        {state.projects.map(p => (
          <Grid item onClick={() => navigateToProject(p)}>
            <ProjectCard project={p} onRemove={() => triggerProjectRemove(p)} />
          </Grid>
        ))}
      </Grid>
    </Box>
  );
};

export default Catalog;
