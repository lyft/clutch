import React from "react";
import { BrowserRouter as Router, Outlet, Route, Routes, useLocation } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
import _ from "lodash";

import AppLayout from "../AppLayout";
import { ApplicationContext } from "../Contexts/app-context";
import type { StorageContextProps } from "../Contexts/storage-context";
import { StorageContext, StorageState } from "../Contexts/storage-context";
import storageContextReducer from "../Contexts/storage-context-reducer";
import { FEATURE_FLAG_POLL_RATE, featureFlags } from "../flags";
import Landing from "../landing";
import { useNavigate } from "../navigation";
import { client } from "../Network";
import type { ClutchError } from "../Network/errors";
import NotFound from "../not-found";

import { registeredWorkflows } from "./registrar";
import { Theme } from "./themes";
import type { ConfiguredRoute, Workflow, WorkflowConfiguration } from "./workflow";
import ErrorBoundary from "./workflow";

export interface UserConfiguration {
  [packageName: string]: {
    [key: string]: ConfiguredRoute;
  };
}

const defaultStorageState: StorageState = {
  shortLinked: false,
  store: {},
  tempStore: {},
};

interface ClutchAppProps {
  availableWorkflows: {
    [key: string]: () => WorkflowConfiguration;
  };
  configuration: UserConfiguration;
}

interface ShortLinkHydratorProps {
  hydrate: (HydrateData) => void;
}

/**
 * Filter workflow routes using available feature flags.
 * @param workflows a list of valid Workflow objects.
 */
const featureFlagFilter = (workflows: Workflow[]): Promise<Workflow[]> => {
  return featureFlags().then(flags =>
    workflows.filter(workflow => {
      /* eslint-disable-next-line no-param-reassign */
      workflow.routes = workflow.routes.filter(route => {
        const show =
          route.featureFlag === undefined ||
          (flags?.[route.featureFlag] !== undefined &&
            flags[route.featureFlag].booleanValue === true);
        return show;
      });
      return workflow.routes.length !== 0;
    })
  );
};

const ShortLinkHydrator = ({ hydrate }: ShortLinkHydratorProps) => {
  const { pathname } = useLocation();
  const navigate = useNavigate();

  React.useEffect(() => {
    // Looking for a route similar to: "/sl/1234"
    const matches = pathname.match(/(\/sl\/)(.*)/i);
    if (matches && matches[2]) {
      const requestData: IClutch.shortlink.v1.IGetRequest = { hash: matches[2] };

      client
        .post("/v1/shortlink/get", requestData)
        .then(response => {
          const { path = "/", state } = response.data as IClutch.shortlink.v1.IGetResponse;
          hydrate(state);
          navigate(path);
        })
        .catch((error: ClutchError) => {
          // eslint-disable-next-line no-console
          console.warn(`Shortlink ${matches[2]} errored, redirecting home`);
          navigate("/");
        });
    }
  }, [pathname]);

  return null;
};

const ClutchApp: React.FC<ClutchAppProps> = ({
  availableWorkflows,
  configuration: userConfiguration,
}) => {
  const [workflows, setWorkflows] = React.useState<Workflow[]>([]);
  const [isLoading, setIsLoading] = React.useState<boolean>(true);

  const [storageState, dispatch] = React.useReducer(storageContextReducer, defaultStorageState);

  const loadWorkflows = () => {
    registeredWorkflows(availableWorkflows, userConfiguration, [featureFlagFilter]).then(w => {
      setWorkflows(w);
      setIsLoading(false);
    });
  };

  const retrieveData = (componentName: string, key: string, defaultData?: any): any => {
    const { store } = storageState;

    // fetch hydrated data
    if (store && store[componentName]) {
      return key.length ? store[componentName][key] : store[componentName];
    }

    // fetch from local storage
    if (key.length) {
      const localData = window.localStorage.getItem(key);
      if (localData) {
        try {
          return JSON.parse(localData);
        } catch (e) {
          return localData;
        }
      }
    }

    return defaultData;
  };

  const storageProviderProps: StorageContextProps = {
    shortLinked: storageState.shortLinked,
    functions: {
      storeData: (componentName: string, key: string, data: any, local?: boolean) =>
        dispatch({ type: "STORE_DATA", payload: { componentName, key, data, local } }),
      removeData: (componentName: string, key: string, local?: boolean) =>
        dispatch({ type: "REMOVE_DATA", payload: { componentName, key, local } }),
      clearData: () => dispatch({ type: "EMPTY_DATA" }),
      tempData: () => storageState.tempStore,
      retrieveData,
    },
  };

  React.useEffect(() => {
    loadWorkflows();
    const interval = setInterval(loadWorkflows, FEATURE_FLAG_POLL_RATE);
    return () => clearInterval(interval);
  }, []);

  const [discoverableWorkflows, setDiscoverableWorkflows] = React.useState([]);
  React.useEffect(() => {
    /** Filter out all of the workflows that are configured to be `hideNav: true`.
     * This prevents the workflows from being discoverable by the user from the UI,
     * both search and drawer navigation.
     *
     * The routes for all configured workflows will still be reachable
     * by manually providing the full path in the URI.
     */
    const pw = _.cloneDeep(workflows).filter(workflow => {
      const publicRoutes = workflow.routes.filter(route => {
        return !(route?.hideNav !== undefined ? route.hideNav : false);
      });
      workflow.routes = publicRoutes; /* eslint-disable-line no-param-reassign */
      return publicRoutes.length !== 0;
    });
    setDiscoverableWorkflows(pw);
  }, [workflows]);

  return (
    <Router>
      <Theme variant="light">
        <div id="App">
          <ApplicationContext.Provider value={{ workflows: discoverableWorkflows }}>
            <StorageContext.Provider value={storageProviderProps}>
              <Routes>
                <Route path="/" element={<AppLayout isLoading={isLoading} />}>
                  <Route key="landing" path="" element={<Landing />} />
                  {workflows.map((workflow: Workflow) => {
                    const workflowPath = workflow.path.replace(/^\/+/, "").replace(/\/+$/, "");
                    const workflowKey = workflow.path.split("/")[0];
                    return (
                      <Route
                        path={`${workflowPath}/`}
                        key={workflowKey}
                        element={
                          <ErrorBoundary workflow={workflow}>
                            <Outlet />
                          </ErrorBoundary>
                        }
                      >
                        {workflow.routes.map(route => {
                          const heading = route.displayName
                            ? `${workflow.displayName}: ${route.displayName}`
                            : workflow.displayName;
                          return (
                            <Route
                              key={workflow.path}
                              path={`${route.path.replace(/^\/+/, "").replace(/\/+$/, "")}`}
                              element={React.cloneElement(<route.component />, {
                                ...route.componentProps,
                                heading,
                              })}
                            />
                          );
                        })}
                        <Route key={`${workflow.path}/notFound`} path="*" element={<NotFound />} />
                      </Route>
                    );
                  })}
                  <Route
                    key="short-links"
                    path="/sl/*"
                    element={
                      <ShortLinkHydrator
                        hydrate={data => dispatch({ type: "HYDRATE", payload: { data } })}
                      />
                    }
                  />
                  <Route key="notFound" path="*" element={<NotFound />} />
                </Route>
              </Routes>
            </StorageContext.Provider>
          </ApplicationContext.Provider>
        </div>
      </Theme>
    </Router>
  );
};

export default ClutchApp;
