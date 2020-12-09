import React from "react";
import { BrowserRouter as Router, Outlet, Route, Routes } from "react-router-dom";
import _ from "lodash";

import AppLayout from "../AppLayout";
import { ApplicationContext } from "../Contexts/app-context";
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

interface ClutchAppProps {
  availableWorkflows: {
    [key: string]: () => WorkflowConfiguration;
  };
  configuration: UserConfiguration;
}

const ClutchApp: React.FC<ClutchAppProps> = ({ availableWorkflows, configuration }) => {
  const workflows = registeredWorkflows(availableWorkflows, configuration);
  const appConfiguration = configuration?.["@clutch-sh/app"] || {};

  /** Filter out all of the workflows that are configured to be `hideNav: true`.
   * This prevents the workflows from being discoverable by the user from the UI,
   * both search and drawer navigation.
   *
   * The routes for all configured workflows will still be reachable
   * by manually providing the full path in the URI.
   */
  const publicWorkflows = _.cloneDeep(workflows).filter(workflow => {
    const publicRoutes = workflow.routes.filter(route => {
      return !(route?.hideNav !== undefined ? route.hideNav : false);
    });
    workflow.routes = publicRoutes; /* eslint-disable-line no-param-reassign */
    return publicRoutes.length !== 0;
  });

  return (
    <Router>
      <Theme variant="light">
        <div id="App">
          <ApplicationContext.Provider
            value={{ workflows: publicWorkflows, appConfig: appConfiguration }}
          >
            <Routes>
              <Route path="/*" element={<AppLayout />}>
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
