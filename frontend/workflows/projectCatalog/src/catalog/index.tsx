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
import RestoreIcon from "@mui/icons-material/Restore";
import SearchIcon from "@mui/icons-material/Search";
import { Box, CircularProgress, Theme } from "@mui/material";

import type { WorkflowProps } from "../types";

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

const PlaceholderContainer = styled("div")(({ theme }: { theme: Theme }) => ({
  margin: theme.spacing(theme.clutch.spacing.lg),
}));

const Placeholder = () => (
  <Paper>
    <PlaceholderContainer>
      <Typography variant="h5">There is nothing to display here</Typography>
      <Typography variant="body3">Please enter a project to proceed.</Typography>
    </PlaceholderContainer>
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
    caseSensitive: false,
  });

  return { results: response?.data?.results || [] };
};

const FormWrapper = styled("div")(({ theme }: { theme: Theme }) => ({
  margin: theme.spacing(theme.clutch.spacing.base),
}));

const Form = styled.form({});

const MainContentWrapper = styled("div")(({ theme }: { theme: Theme }) => ({
  display: "flex",
  justifyContent: "space-between",
  marginBottom: theme.spacing(theme.clutch.spacing.base),
  marginTop: theme.spacing(theme.clutch.spacing.lg),
}));

const PlaceholderWrapper = styled(Grid)(({ theme }: { theme: Theme }) => ({
  paddingTop: theme.spacing(theme.clutch.spacing.lg),
}));

const Catalog: React.FC<WorkflowProps> = ({ allowDisabled }) => {
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
      !hasState(),
      allowDisabled
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
      },
      allowDisabled
    );
  };

  const triggerProjectRemove = project => {
    removeProject(
      project.name,
      projects => dispatch({ type: "REMOVE_PROJECT", payload: { projects } }),
      setError,
      allowDisabled
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
    <Box>
      <Paper>
        <FormWrapper>
          <Form noValidate onSubmit={handleSubmit(triggerProjectAdd)}>
            <TextField
              label="Search"
              placeholder="Project Name"
              value={state.search}
              onChange={handleChanges}
              autocompleteCallback={v => autoComplete(v)}
              endAdornment={state.isSearching ? <CircularProgress size="24px" /> : <SearchIcon />}
              error={state.error !== undefined}
              helperText={state?.error}
            />
          </Form>
        </FormWrapper>
      </Paper>
      <MainContentWrapper>
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
                true,
                allowDisabled
              );
            }}
          >
            <RestoreIcon />
          </IconButton>
        </Tooltip>
      </MainContentWrapper>
      {state.projects.length ? (
        <Grid container direction="row" spacing={3}>
          {state.projects.map(p => (
            <Grid item key={p.name} onClick={() => navigateToProject(p)}>
              <ProjectCard project={p} onRemove={() => triggerProjectRemove(p)} />
            </Grid>
          ))}
        </Grid>
      ) : (
        <PlaceholderWrapper container justifyContent="center">
          <Placeholder />
        </PlaceholderWrapper>
      )}
    </Box>
  );
};

export default Catalog;
