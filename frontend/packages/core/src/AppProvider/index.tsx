import React from "react";
import { BrowserRouter as Router, Outlet, Route, Routes } from "react-router-dom";
import Bugsnag from "@bugsnag/js";
import BugsnagPluginReact from "@bugsnag/plugin-react";

import AppLayout from "../AppLayout";
import AppNotification from "../AppNotifications";
import { ApplicationContext, ShortLinkContext, UserPreferencesProvider } from "../Contexts";
import type { HeaderItem, TriggeredHeaderData } from "../Contexts/app-context";
import type { ShortLinkContextProps } from "../Contexts/short-link-context";
import type { HydratedData, HydrateState } from "../Contexts/workflow-storage-context/types";
import { Toast } from "../Feedback";
import { FEATURE_FLAG_POLL_RATE, featureFlags } from "../flags";
import Landing from "../landing";
import type { ClutchError } from "../Network/errors";
import NotFound from "../not-found";
import type { AppConfiguration } from "../Types";
import WorkflowLayout, { LayoutProps } from "../WorkflowLayout";

import { registeredWorkflows } from "./registrar";
import ShortLinkProxy, { ShortLinkBaseRoute } from "./short-link-proxy";
import ShortLinkStateHydrator from "./short-link-state-hydrator";
import { Theme } from "./themes";
import type { ConfiguredRoute, Workflow, WorkflowConfiguration } from "./workflow";
import ErrorBoundary from "./workflow";

export interface WorkflowIcon {
  path: string;
}

export interface UserConfiguration {
  [packageName: string]: {
    icon: WorkflowIcon;
    [key: string]: WorkflowIcon | ConfiguredRoute;
  };
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

enum ChildType {
  HEADER = "header",
}

interface ChildTypes {
  type: ChildType;
}

type ClutchAppChildType = ChildTypes;

type ClutchAppChild = React.ReactElement<ClutchAppChildType>;

interface ClutchAppProps {
  availableWorkflows: {
    [key: string]: () => WorkflowConfiguration;
  };
  configuration: UserConfiguration;
  appConfiguration?: AppConfiguration;
  children?: ClutchAppChild | ClutchAppChild[];
}

const ClutchApp = ({
  availableWorkflows,
  configuration: userConfiguration,
  appConfiguration,
  children,
}: ClutchAppProps) => {
  const [workflows, setWorkflows] = React.useState<Workflow[]>([]);
  const [isLoading, setIsLoading] = React.useState<boolean>(true);
  const [workflowSessionStore, setWorkflowSessionStore] = React.useState<HydratedData>();
  const [hydrateState, setHydrateState] = React.useState<HydrateState | null>(null);
  const [hydrateError, setHydrateError] = React.useState<ClutchError | null>(null);
  const [hasCustomLanding, setHasCustomLanding] = React.useState<boolean>(false);
  const [triggeredHeaderData, setTriggeredHeaderData] = React.useState<TriggeredHeaderData>();
  const [customHeader, setCustomHeader] = React.useState<React.ReactElement>();

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
          [item]: data as any,
        }),
      triggeredHeaderData,
    }),
    [discoverableWorkflows, triggeredHeaderData]
  );

  React.useEffect(() => {
    if (children) {
      React.Children.forEach(children, child => {
        if (React.isValidElement(child)) {
          const { type = "" } = child?.props || {};

          if (type.toLowerCase() === ChildType.HEADER) {
            setCustomHeader(child);
          }
        }
      });
    }
  }, [children]);

  return (
    <Router>
      <Theme>
        <div id="App">
          <ApplicationContext.Provider value={appContextValue}>
            <UserPreferencesProvider>
              <ShortLinkContext.Provider value={shortLinkProviderProps}>
                {hydrateError && (
                  <Toast onClose={() => setHydrateError(null)}>
                    Unable to retrieve short link: {hydrateError?.status?.text}
                  </Toast>
                )}
                <Routes>
                  <Route
                    path="/"
                    element={
                      <AppLayout
                        isLoading={isLoading}
                        configuration={appConfiguration}
                        header={customHeader}
                      />
                    }
                  >
                    {!hasCustomLanding && (
                      <Route
                        key="landing"
                        path=""
                        element={
                          <AppNotification
                            type="layout"
                            workflow="landing"
                            banners={appConfiguration?.banners}
                          >
                            <Landing />
                          </AppNotification>
                        }
                      />
                    )}
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

                            // We define these props in order to avoid UI changes before refactoring
                            const workflowLayoutProps: LayoutProps = {
                              ...route.layoutProps,
                              heading: route.layoutProps?.heading || heading,
                              workflow,
                            };

                            const workflowRouteComponent = (
                              <AppNotification
                                type="layout"
                                workflow={workflow?.displayName}
                                banners={appConfiguration?.banners}
                              >
                                {React.cloneElement(<route.component />, {
                                  ...route.componentProps,
                                  // This is going to be removed to be used in the WorkflowLayout only
                                  heading,
                                })}
                              </AppNotification>
                            );

                            return (
                              <Route
                                key={workflow.path}
                                path={`${route.path.replace(/^\/+/, "").replace(/\/+$/, "")}`}
                                element={
                                  appConfiguration?.enableWorkflowLayout ? (
                                    <WorkflowLayout {...workflowLayoutProps}>
                                      {workflowRouteComponent}
                                    </WorkflowLayout>
                                  ) : (
                                    workflowRouteComponent
                                  )
                                }
                              />
                            );
                          })}
                          <Route
                            key={`${workflow.path}/notFound`}
                            path="*"
                            element={<NotFound />}
                          />
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
            </UserPreferencesProvider>
          </ApplicationContext.Provider>
        </div>
      </Theme>
    </Router>
  );
};

const BugSnagApp = (props: ClutchAppProps) => {
  if (process.env.REACT_APP_BUGSNAG_API_TOKEN) {
    // eslint-disable-next-line no-underscore-dangle
    if (!(Bugsnag as any)._client) {
      Bugsnag.start({
        apiKey: process.env.REACT_APP_BUGSNAG_API_TOKEN,
        plugins: [new BugsnagPluginReact()],
        releaseStage: process.env.APPLICATION_ENV || "production",
      });
    }
    const BugsnagBoundary = Bugsnag.getPlugin("react").createErrorBoundary(React);
    return (
      <BugsnagBoundary>
        <ClutchApp {...props} />
      </BugsnagBoundary>
    );
  }

  return <ClutchApp {...props} />;
};

export default BugSnagApp;
