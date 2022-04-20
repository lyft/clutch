import React from "react";
import { useForm } from "react-hook-form";
import {
  client,
  Grid,
  IconButton,
  Paper,
  TextField,
  Tooltip,
  Typography,
  useNavigate,
} from "@clutch-sh/core";
import styled from "@emotion/styled";
import { Box, CircularProgress } from "@material-ui/core";
import RestoreIcon from "@material-ui/icons/Restore";
import SearchIcon from "@material-ui/icons/Search";

import type { WorkflowProps } from "..";

import catalogReducer from "./catalog-reducer";
import ProjectCard from "./project-card";
import { addProject, clearProjects, getProjects, hasState, removeProject } from "./storage";
import type { CatalogState } from "./types";

const initialState: CatalogState = {
  projects: [],
  search: "",
  isLoading: false,
  isSearching: false,
  error: undefined,
};

const Placeholder = () => (
  <Paper>
    <div style={{ margin: "32px", textAlign: "center" }}>
      <Typography variant="h5">There is nothing to display here</Typography>
      <Typography variant="body3">Please enter a project to proceed.</Typography>
    </div>
  </Paper>
);

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

const Form = styled.form({});

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
      setError,
      !hasState()
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
        const projectMatches = projects.filter(
          p => state?.search && state.search !== "" && p?.name === state.search
        );
        if (projectMatches.length === 1) {
          navigateToProject(projectMatches[0]);
        }
      },
      e => {
        dispatch({ type: "SEARCH_END" });
        setError(e);
      }
    );
  };

  const triggerProjectRemove = project => {
    removeProject(
      project.name,
      projects => dispatch({ type: "REMOVE_PROJECT", payload: { projects } }),
      setError
    );
  };

  const { handleSubmit } = useForm({
    mode: "onSubmit",
    reValidateMode: "onSubmit",
    shouldFocusError: false,
  });

  const handleChanges = event => {
    dispatch({ type: "SEARCH", payload: { search: event.target.value } });
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
        <div style={{ marginTop: "8px" }}>
          <Typography variant="subtitle3" color="rgb(13, 16, 48, .48)">
            A catalog of all projects.
          </Typography>
        </div>
      </div>
      <Paper>
        <div style={{ margin: "16px" }}>
          <Form noValidate onSubmit={handleSubmit(triggerProjectAdd)}>
            <TextField
              placeholder="Search"
              value={state.search}
              onChange={handleChanges}
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
          </Form>
        </div>
      </Paper>
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          marginBottom: "16px",
          marginTop: "32px",
        }}
      >
        <Typography variant="h3">My Projects</Typography>
        <Tooltip title="Restore to owned projects only">
          <IconButton
            variant="neutral"
            onClick={() => {
              clearProjects();
              dispatch({ type: "HYDRATE_START" });
              getProjects(
                projects => dispatch({ type: "HYDRATE_END", payload: { result: projects } }),
                setError,
                true
              );
            }}
          >
            <RestoreIcon />
          </IconButton>
        </Tooltip>
      </div>
      {state.projects.length ? (
        <Grid container direction="row" spacing={3}>
          {state.projects.map(p => (
            <Grid item onClick={() => navigateToProject(p)}>
              <ProjectCard project={p} onRemove={() => triggerProjectRemove(p)} />
            </Grid>
          ))}
        </Grid>
      ) : (
        <Grid container justify="center" style={{ paddingTop: "35px" }}>
          <Placeholder />
        </Grid>
      )}
    </Box>
  );
};

export default Catalog;
