import React from "react";
import { BrowserRouter as Router, Outlet, Route, Routes } from "react-router-dom";
import { CssBaseline, MuiThemeProvider } from "@material-ui/core";
import { ThemeProvider as StyledThemeProvider } from "styled-components";

import AppLayout from "../AppLayout";
import { ApplicationContext } from "../Contexts/app-context";
import Landing from "../landing";
import NotFound from "../not-found";

import registeredWorkflows from "./registrar";
import { getTheme } from "./themes";
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
  return (
    <Router>
      <MuiThemeProvider theme={getTheme()}>
        <StyledThemeProvider theme={getTheme()}>
          <CssBaseline />
          <div id="App">
            <ApplicationContext.Provider value={{ workflows }}>
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
        </StyledThemeProvider>
      </MuiThemeProvider>
    </Router>
  );
};

export default ClutchApp;
