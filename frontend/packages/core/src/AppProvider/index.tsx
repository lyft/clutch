import React from "react";
import { BrowserRouter as Router, Outlet, Route, Routes } from "react-router-dom";

import AppLayout from "../AppLayout";
import { ApplicationContext, ShortLinkContext } from "../Contexts";
import type { HeaderItem, TriggeredHeaderData } from "../Contexts/app-context";
import type { ShortLinkContextProps } from "../Contexts/short-link-context";
import type { HydratedData, HydrateState } from "../Contexts/workflow-storage-context/types";
import { Toast } from "../Feedback";
import { FEATURE_FLAG_POLL_RATE, featureFlags } from "../flags";
import Landing from "../landing";
import type { ClutchError } from "../Network/errors";
import NotFound from "../not-found";

import { registeredWorkflows } from "./registrar";
import ShortLinkProxy, { ShortLinkBaseRoute } from "./short-link-proxy";
import ShortLinkStateHydrator from "./short-link-state-hydrator";
import { Theme } from "./themes";
import type { ConfiguredRoute, Workflow, WorkflowConfiguration } from "./workflow";
import ErrorBoundary from "./workflow";

export interface UserConfiguration {
  [packageName: string]: {
    [key: string]: ConfiguredRoute;
  };
}

export interface AppConfiguration {
  /** Will override the title of the given application */
  title?: string;
  /** Supports a react node or a string representing a public assets path */
  logo?: React.ReactNode | string;
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

interface ClutchAppProps {
  availableWorkflows: {
    [key: string]: () => WorkflowConfiguration;
  };
  configuration: UserConfiguration;
  appConfiguration: AppConfiguration;
}

const ClutchApp: React.FC<ClutchAppProps> = ({
  availableWorkflows,
  configuration: userConfiguration,
  appConfiguration,
}) => {
  const [workflows, setWorkflows] = React.useState<Workflow[]>([]);
  const [isLoading, setIsLoading] = React.useState<boolean>(true);
  const [workflowSessionStore, setWorkflowSessionStore] = React.useState<HydratedData>();
  const [hydrateState, setHydrateState] = React.useState<HydrateState | null>(null);
  const [hydrateError, setHydrateError] = React.useState<ClutchError | null>(null);
  const [hasCustomLanding, setHasCustomLanding] = React.useState<boolean>(false);
  const [triggeredHeaderData, setTriggeredHeaderData] = React.useState<TriggeredHeaderData>();

  /** Used to control a race condition from displaying the workflow and the state being updated with the hydrated data */
  const [shortLinkLoading, setShortLinkLoading] = React.useState<boolean>(false);

  const loadWorkflows = () => {
    registeredWorkflows(availableWorkflows, userConfiguration, [featureFlagFilter]).then(w => {
      setWorkflows(w);
      setIsLoading(false);
    });
  };

  React.useEffect(() => {
    loadWorkflows();
    const interval = setInterval(loadWorkflows, FEATURE_FLAG_POLL_RATE);
    return () => clearInterval(interval);
  }, []);

  const [discoverableWorkflows, setDiscoverableWorkflows] = React.useState([]);
  React.useEffect(() => {
    const landingWorkflows = workflows.filter(workflow => workflow.path === "");
    if (landingWorkflows.length > 0) {
      /** Used to control a custom landing page */
      setHasCustomLanding(true);
    }
    setDiscoverableWorkflows(workflows);
  }, [workflows]);

  const shortLinkProviderProps: ShortLinkContextProps = React.useMemo(
    () => ({
      removeWorkflowSession: () => setWorkflowSessionStore(null),
      retrieveWorkflowSession: () => workflowSessionStore,
      storeWorkflowSession: setWorkflowSessionStore,
    }),
    [workflowSessionStore]
  );

  const appContextValue = React.useMemo(
    () => ({
      workflows: discoverableWorkflows,
      triggerHeaderItem: (item: HeaderItem, data: unknown) =>
        // Will set the open status and spread any additional data onto the property named after the item
        setTriggeredHeaderData({
          ...triggeredHeaderData,
          [item]: {
            ...(data as any),
          },
        }),
      triggeredHeaderData,
    }),
    [discoverableWorkflows, triggeredHeaderData]
  );

  return (
    <Router>
      {/* TODO: use the ThemeProvider for proper theming in the future 
        See https://github.com/lyft/clutch/commit/f6c6706b9ba29c4d4c3e5d0ac0c5d0f038203937 */}
      <Theme variant="light">
        <div id="App">
          <ApplicationContext.Provider value={appContextValue}>
            <ShortLinkContext.Provider value={shortLinkProviderProps}>
              {hydrateError && (
                <Toast onClose={() => setHydrateError(null)}>
                  Unable to retrieve short link: {hydrateError?.status?.text}
                </Toast>
              )}
              <Routes>
                <Route
                  path="/"
                  element={<AppLayout isLoading={isLoading} configuration={appConfiguration} />}
                >
                  {!hasCustomLanding && <Route key="landing" path="" element={<Landing />} />}
                  {workflows.map((workflow: Workflow) => {
                    const workflowPath = workflow.path.replace(/^\/+/, "").replace(/\/+$/, "");
                    const workflowKey = workflow.path.split("/")[0];
                    return (
                      <Route
                        path={`${workflowPath}/`}
                        key={workflowKey}
                        element={
                          <ErrorBoundary workflow={workflow}>
                            <ShortLinkStateHydrator
                              sharedState={hydrateState}
                              clearTemporaryState={() => setHydrateState(null)}
                            >
                              {!shortLinkLoading && <Outlet />}
                            </ShortLinkStateHydrator>
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
                    path={`/${ShortLinkBaseRoute}/:hash`}
                    element={
                      <ShortLinkProxy
                        setLoading={setShortLinkLoading}
                        hydrate={setHydrateState}
                        onError={setHydrateError}
                      />
                    }
                  />
                  <Route key="notFound" path="*" element={<NotFound />} />
                </Route>
              </Routes>
            </ShortLinkContext.Provider>
          </ApplicationContext.Provider>
        </div>
      </Theme>
    </Router>
  );
};

export default ClutchApp;
