import * as React from "react";
import { useForm } from "react-hook-form";
import type { clutch as IClutch, google as IGoogle } from "@clutch-sh/api";
import {
  client,
  ClutchError,
  FeatureOn,
  SimpleFeatureFlag,
  styled,
  TextField,
  Tooltip,
  TooltipContainer,
  Typography,
  userId,
  useTheme,
  useWorkflowStorageContext,
} from "@clutch-sh/core";
import AddIcon from "@mui/icons-material/Add";
import InfoOutlinedIcon from "@mui/icons-material/InfoOutlined";
import UpdateIcon from "@mui/icons-material/Update";
import { alpha, Divider, LinearProgress, Theme } from "@mui/material";
import _ from "lodash";

import { useDashUpdater, useRefreshRateState, useRefreshUpdater } from "./dash-hooks";
import {
  deriveStateData,
  DispatchContext,
  exclusiveProjectDependencies,
  StateContext,
} from "./helpers";
import ProjectGroup from "./project-group";
import selectorReducer from "./selector-reducer";
import { COMPONENT_NAME, getLocalState, loadStoredState, STORAGE_STATE_KEY } from "./storage";
import type { Action, DashState, State } from "./types";
import { Group } from "./types";

/**
 * ProjectSelectorError: use for defining error types that will be propagated to the parent component
 */
export interface ProjectSelectorError {
  /**
   * errors: Partial failure errors from the server
   */
  errors: IGoogle.rpc.IStatus[];
}

/**
 * ProjectSelectorProps: Defined input properties of ProjectSelector component
 */
interface ProjectSelectorProps {
  /**
   * onError: (optional) error handle which will accept a ProjectSelectorError as input
   */
  onError?: (ProjectSelectorError) => void;
}

const initialState: State = {
  [Group.PROJECTS]: {},
  [Group.UPSTREAM]: {},
  [Group.DOWNSTREAM]: {},
  projectData: {},
  projectErrors: [],
  loading: false,
  error: undefined,
};

const StyledSelectorContainer = styled("div")(({ theme }: { theme: Theme }) => ({
  backgroundColor: theme.palette.primary[50],
  borderRight: `1px solid ${alpha(theme.palette.secondary[900], 0.1)}`,
  boxShadow: `0px 4px 6px ${alpha(theme.palette.primary[600], 0.2)}`,
  width: "245px",
  overflowY: "auto",
  overflowX: "hidden",
  maxHeight: "100%",
}));

const StyledWorkflowHeader = styled("div")({
  margin: "16px 16px 12px 16px",
  display: "flex",
  alignItems: "center",
  justifyContent: "space-between",
  height: "24px",
});

const StyledWorkflowTitle = styled("span")({
  fontWeight: "bold",
  fontSize: "20px",
  lineHeight: "24px",
  margin: "0px 8px",
});

const StyledProjectTextField = styled(TextField)({
  padding: "16px 16px 8px 16px",
});

const StyledProgressContainer = styled("div")(({ theme }: { theme: Theme }) => ({
  height: "4px",
  ".MuiLinearProgress-root": {
    backgroundColor: theme.palette.primary[400],
  },
  ".MuiLinearProgress-bar": {
    backgroundColor: theme.palette.primary[600],
  },
}));

// TODO(smonero): decide on styling for this
const FlexCenterAlignContainer = styled("div")({
  display: "flex",
  alignItems: "center",
});

// Determines if every project has projectData (i.e. the effect has finished fetching the data)
const allPresent = (state: State, dispatch: React.Dispatch<Action>): boolean => {
  const projectCheck = project =>
    project in state.projectData && !_.isEmpty(state.projectData?.[project]);

  // verify that we have all data for our top level projects
  const projectKeys = Object.keys(state[Group.PROJECTS]);
  if (!projectKeys.every(projectCheck)) {
    // missing data for a direct project, return
    return false;
  }

  // we have data for our projects, pull the dependencies from each to check
  // (instead of relying on our state)
  const dependencies: string[] = [];
  projectKeys.forEach(p => {
    const { upstreams, downstreams } = exclusiveProjectDependencies(state, Group.PROJECTS, p);
    dependencies.push(...upstreams, ...downstreams);
  });

  if (!_.uniq(dependencies).every(projectCheck)) {
    // missing data for a dependency, return
    return false;
  }

  const upDownKeys: string[] = Array.from(
    new Set([...Object.keys(state[Group.UPSTREAM]), ...Object.keys(state[Group.DOWNSTREAM])])
  );

  // we've received all data for our projects, check for mismatch in up / down streams in our state and remove them
  const missing = upDownKeys.filter(p => !projectCheck(p));
  if (missing.length) {
    // we're missing upstreams / downstreams, remove them from our state
    dispatch({ type: "HYDRATE_ERROR", payload: { result: { missing } } });
  }

  return true;
};

const hydrateProjects = (
  state: State,
  dispatch: React.Dispatch<Action>,
  fromShortLink: boolean
) => {
  // Determine if any hydration is required.
  // - Are any services missing from state.projectdata?
  // - Are projects empty (first load)?
  // - Is loading not already in progress?
  if (
    !state.loading &&
    (Object.keys(state[Group.PROJECTS]).length === 0 || !allPresent(state, dispatch))
  ) {
    dispatch({ type: "HYDRATE_START" });

    // TODO: have userId check be server driven
    const requestParams = { users: [], projects: [] } as {
      users: string[];
      projects: string[];
    };

    // will only push the userId onto the request if we're not currently under a shortLink
    // this is so that the API will not return the users default projects
    if (!fromShortLink) {
      requestParams.users.push(userId());
    }

    // since default projects are always included in the response, no reason to only filter on custom projects
    requestParams.projects = Object.keys(state[Group.PROJECTS]);

    client
      .post("/v1/project/getProjects", requestParams as IClutch.project.v1.GetProjectsRequest)
      .then(resp => {
        const { results, partialFailures } = resp.data as IClutch.project.v1.GetProjectsResponse;

        dispatch({
          type: "HYDRATE_END",
          payload: { result: { results: results || {}, partialFailures } },
        });
      })
      .catch((err: ClutchError) => {
        dispatch({ type: "HYDRATE_ERROR", payload: { result: err } });
      });
  }
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
    caseSensitive: false,
  });

  return { results: response?.data?.results || [] };
};

const Form = styled("form")({});

const ProjectSelector = ({ onError }: ProjectSelectorProps) => {
  const theme = useTheme();
  // On load, we'll request a list of owned projects and their upstreams and downstreams from the API.
  // The API will contain information about the relationships between projects and upstreams and downstreams.
  // By default, the owned projects will be checked and others will be unchecked.

  const [customProject, setCustomProject] = React.useState("");
  const { updateSelected } = useDashUpdater();

  const { removeData, retrieveData, fromShortLink, storeData } = useWorkflowStorageContext();

  const [state, dispatch] = React.useReducer(
    selectorReducer,
    loadStoredState(initialState, retrieveData, removeData)
  );

  React.useEffect(() => {
    hydrateProjects(state, dispatch, fromShortLink);
  }, [state[Group.PROJECTS]]);

  // computes the final state for rendering across other components
  // (ie. filters out upstream/downstreams that are "hidden")
  const derivedState = React.useMemo(() => deriveStateData(state), [state]);

  // This hook updates the global dash state based on the currently selected projects for cards to consume (including upstreams and downstreams).
  React.useEffect(() => {
    if (!allPresent(state, dispatch)) {
      // Need to wait for the data.
      return;
    }

    if (onError && state.projectErrors && state.projectErrors.length) {
      onError({ errors: state.projectErrors });
      dispatch({ type: "CLEAR_PROJECT_ERRORS" });
    }

    const dashState: DashState = { projectData: {}, selected: [] };

    // Determine selected projects.
    const selected = new Set<string>();
    _.forEach(Object.keys(derivedState[Group.PROJECTS]), p => {
      if (derivedState[Group.PROJECTS][p].checked) {
        selected.add(p);
      }
    });
    _.forEach(Object.keys(derivedState[Group.DOWNSTREAM]), p => {
      if (derivedState[Group.DOWNSTREAM][p].checked) {
        selected.add(p);
      }
    });
    _.forEach(Object.keys(derivedState[Group.UPSTREAM]), p => {
      if (derivedState[Group.UPSTREAM][p].checked) {
        selected.add(p);
      }
    });
    dashState.selected = Array.from(selected).sort();

    // Collect project data.
    _.forEach(dashState.selected, p => {
      dashState.projectData[p] = state.projectData[p];
    });

    // Update!
    storeData(COMPONENT_NAME, STORAGE_STATE_KEY, getLocalState(state), true);
    updateSelected(dashState);
  }, [state]);

  const handleAdd = v => {
    if (!v.project || v?.project === "") {
      return;
    }
    dispatch({
      type: "ADD_PROJECTS",
      payload: { group: Group.PROJECTS, projects: [v.project] },
    });
    setCustomProject("");
  };

  const hasError = state.error !== undefined && state.error !== null;

  const { handleSubmit, register } = useForm({
    mode: "onSubmit",
    reValidateMode: "onSubmit",
    shouldFocusError: false,
  });

  const handleChanges = event => {
    setCustomProject(event.target.value);
  };

  const handleSubmission = _.debounce(() => {
    handleSubmit(handleAdd)();
  }, 100);

  const { updateRefreshRate } = useRefreshUpdater();

  const initialRefreshRateState = useRefreshRateState();
  const { refreshRate } = initialRefreshRateState;
  // if refreshRate is null, that means autoRefresh is disabled
  // else if its a number, it will be used as the value for setInterval
  const [autoRefresh, setAutoRefresh] = React.useState<boolean>(!!refreshRate);

  const handleRefreshRateChange = () => {
    // if autorefresh was false, it is now true, so set it to the value
    if (!autoRefresh) {
      updateRefreshRate({ refreshRate: 30000 });
    } else {
      // else, it was true, so now it is false, and we should turn it off
      updateRefreshRate({ refreshRate: null });
    }
    setAutoRefresh(!autoRefresh);
  };

  // useEffect that handles auto-refreshing / reloading.
  React.useEffect(() => {
    // This way of autorefreshing is also used in other places.
    if (refreshRate) {
      const interval = setInterval(
        () => hydrateProjects(state, dispatch, fromShortLink),
        refreshRate
      );
      return () => clearInterval(interval);
    }
    return () => {};
  }, [state, refreshRate]);

  return (
    <DispatchContext.Provider value={dispatch}>
      <StateContext.Provider value={derivedState}>
        <StyledSelectorContainer>
          <StyledWorkflowHeader>
            <FlexCenterAlignContainer>
              <StyledWorkflowTitle>Dash</StyledWorkflowTitle>
              <Tooltip
                title={
                  <>
                    {[
                      {
                        title: "Projects",
                        description:
                          "Service, mobile app, etc. Unchecking a project hides its upstream and downstream dependencies.",
                      },
                      {
                        title: "Upstreams",
                        description: "Receive requests and send responses to the selected project.",
                      },
                      {
                        title: "Downstreams",
                        description:
                          "Send requests and receive responses from the selected project.",
                      },
                    ].map(item => (
                      <TooltipContainer key={item.title}>
                        <Typography variant="subtitle3" color={theme.palette.contrastColor}>
                          {item.title}
                        </Typography>
                        <Typography variant="body4" color={theme.palette.secondary[200]}>
                          {item.description}
                        </Typography>
                      </TooltipContainer>
                    ))}
                  </>
                }
                maxWidth="400px"
                placement="right-start"
              >
                <InfoOutlinedIcon fontSize="small" />
              </Tooltip>
            </FlexCenterAlignContainer>
            <FlexCenterAlignContainer>
              <SimpleFeatureFlag feature="dashRefreshSelect">
                <FeatureOn>
                  <Tooltip
                    title={autoRefresh ? "Disable data refresh" : "Data refresh every 30 seconds"}
                    placement="bottom"
                  >
                    <UpdateIcon
                      style={{
                        color: autoRefresh
                          ? theme.palette.primary[600]
                          : alpha(theme.palette.getContrastText(theme.palette.contrastColor), 0.26),
                      }}
                      fontSize="small"
                      onClick={() => {
                        handleRefreshRateChange();
                      }}
                    />
                  </Tooltip>
                </FeatureOn>
              </SimpleFeatureFlag>
            </FlexCenterAlignContainer>
          </StyledWorkflowHeader>
          <StyledProgressContainer>
            {state.loading && <LinearProgress color="secondary" />}
          </StyledProgressContainer>
          <Divider />
          <Form
            noValidate
            onSubmit={e => {
              e.preventDefault();
              handleSubmission();
            }}
          >
            <StyledProjectTextField
              disabled={state.loading}
              placeholder="Add a project"
              name="project"
              value={customProject}
              onChange={handleChanges}
              onReturn={handleSubmission}
              helperText={state.error?.message}
              error={hasError}
              autocompleteCallback={v => autoComplete(v)}
              formRegistration={register}
              endAdornment={<AddIcon />}
            />
          </Form>
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
