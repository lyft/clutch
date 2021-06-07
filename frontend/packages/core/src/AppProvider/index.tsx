import React from "react";
import { BrowserRouter as Router, Outlet, Route, Routes } from "react-router-dom";
import _ from "lodash";

import AppLayout from "../AppLayout";
import { ApplicationContext } from "../Contexts/app-context";
import { FEATURE_FLAG_POLL_RATE, featureFlags } from "../flags";
import Landing from "../landing";
import NotFound from "../not-found";

import registeredWorkflows from "./registrar";
import { Theme } from "./themes";
import type { ConfiguredRoute, Workflow, WorkflowConfiguration } from "./workflow";
import ErrorBoundary from "./workflow";

export interface UserConfiguration {
  [packageName: string]: {
    [key: string]: ConfiguredRoute;
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

interface ClutchAppProps {
  availableWorkflows: {
    [key: string]: () => WorkflowConfiguration;
  };
  configuration: UserConfiguration;
}

const ClutchApp: React.FC<ClutchAppProps> = ({
  availableWorkflows,
  configuration: userConfiguration,
}) => {
  const [workflows, setWorkflows] = React.useState<Workflow[]>([]);
  const [isLoading, setIsLoading] = React.useState<boolean>(true);

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
            <Routes>
              <Route path="/*" element={<AppLayout isLoading={isLoading} />}>
                <Route key="landing" path="/" element={<Landing />} />
                {workflows.map((workflow: Workflow) => (
                  <ErrorBoundary workflow={workflow} key={workflow.path.split("/")[0]}>
                    <Route path={`${workflow.path}/*`} element={<Outlet />}>
                      {workflow.routes.map(route => {
                        const heading = route.displayName
                          ? `${workflow.displayName}: ${route.displayName}`
                          : workflow.displayName;
                        return (
                          <Route
                            key={workflow.path}
                            path={`${route.path}`}
                            element={React.cloneElement(<route.component />, {
                              ...route.componentProps,
                              heading,
                            })}
                          />
                        );
                      })}
                      <Route key={`${workflow.path}/notFound`} path="*" element={<NotFound />} />
                    </Route>
                  </ErrorBoundary>
                ))}
                <Route key="notFound" path="*" element={<NotFound />} />
              </Route>
            </Routes>
          </ApplicationContext.Provider>
        </div>
      </Theme>
    </Router>
  );
};

export default ClutchApp;
