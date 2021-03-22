import React from "react";
import { Grid, IconButton } from "@material-ui/core";
import LaunchIcon from "@material-ui/icons/Launch";
import { Alert } from "@material-ui/lab";

import { Dialog, DialogContent } from "../dialog";
import Code from "../text";

export interface BaseWorkflowProps {
  heading: string;
}

interface Developer {
  contactUrl: string;
  name: string;
}

interface BaseWorkflowConfiguration {
  developer: Developer;
  displayName: string;
  group: string;
  path: string;
  routes: unknown;
}

export interface Workflow extends BaseWorkflowConfiguration {
  routes: ConfiguredRoute[];
}

export interface WorkflowConfiguration extends BaseWorkflowConfiguration {
  routes: {
    [key: string]: Route;
  };
}

interface Route {
  component: React.FC<any>;
  description: string;
  displayName?: string;
  path: string;
  /** Properties required by the Component that are set only via the config. */
  requiredConfigProps?: string[];
  /** Is the workflow discoverable via search and drawer navigation. This defaults to false. */
  hideNav?: boolean;
}

export interface ConfiguredRoute extends Route {
  componentProps?: object;
  trending?: boolean;
}

interface ErrorBoundaryProps {
  workflow: Workflow;
}

interface ErrorBoundaryState {
  error: Error;
  errorInfo: React.ErrorInfo;
  showDetails: boolean;
}

class ErrorBoundary extends React.Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: { workflow: Workflow }) {
    super(props);
    this.state = { error: null, errorInfo: null, showDetails: false };
    this.onDetailsClose = this.onDetailsClose.bind(this);
    this.onDetailsOpen = this.onDetailsOpen.bind(this);
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    this.setState({ error, errorInfo });
  }

  onDetailsClose() {
    this.setState(state => {
      return { ...state, showDetails: false };
    });
  }

  onDetailsOpen() {
    this.setState(state => {
      return { ...state, showDetails: true };
    });
  }

  render() {
    const { children, workflow } = this.props;
    const { error, errorInfo, showDetails } = this.state;

    const defaultErrorMsg = (
      <>Failed to load{workflow.displayName} workflow. Please contact the developer</>
    );
    if (error) {
      let message = defaultErrorMsg;
      if (workflow?.developer) {
        const developerName = ` ${workflow.developer?.name}` || "Unknown";
        const link = (
          <a rel="noopener noreferrer" target="_blank" href={workflow.developer.contactUrl}>
            {developerName}
          </a>
        );
        message = (
          <>
            {defaultErrorMsg}:{link}.
          </>
        );
      }

      return (
        <Grid container direction="column" justify="center" alignItems="center">
          <Dialog onClose={this.onDetailsClose} open={showDetails} title="Stack Trace">
            <DialogContent>
              <Code>{errorInfo?.componentStack || "Could not determine stack trace"}</Code>
            </DialogContent>
          </Dialog>
          <Alert
            severity="error"
            action={
              <IconButton
                aria-label="error"
                color="inherit"
                size="small"
                onClick={this.onDetailsOpen}
              >
                <LaunchIcon />
              </IconButton>
            }
          >
            <div>{message}</div>
          </Alert>
        </Grid>
      );
    }
    return children;
  }
}

export default ErrorBoundary;
