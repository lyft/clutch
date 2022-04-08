import React from "react";
import {
  client,
  Grid,
  IconButton,
  Paper,
  TextField,
  Toast,
  Tooltip,
  Typography,
  useNavigate,
} from "@clutch-sh/core";
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
      <Typography variant="body3">Please enter a namespace and clientset to proceed.</Typography>
    </div>
  </Paper>
);

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
          />
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
        <Grid container direction="row" spacing={5}>
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
      {state.error && (
        <Toast severity="error" onClose={() => dispatch({ type: "CLEAR_ERROR" })}>
          {state.error}
        </Toast>
      )}
    </Box>
  );
};

export default Catalog;
