import React from "react";
import LaunchIcon from "@mui/icons-material/Launch";
import { Alert, Grid, IconButton } from "@mui/material";

import { Dialog, DialogContent } from "../dialog";
import Code from "../text";
import type { LayoutProps } from "../WorkflowLayout";

import type { WorkflowIcon } from "./index";

export interface BaseWorkflowProps {
  heading: string;
}

interface Developer {
  contactUrl: string;
  name: string;
}

interface BaseWorkflowConfiguration {
  /**
   * Team name and contact url, which will be displayed in the case of errors to
   * allow contact of the owning developer.
   */
  developer: Developer;
  /**
   * The name of the workflow which will be displayed in the UI.
   */
  displayName: string;
  /**
   * The name of the group which will display in the sidebar.
   */
  group: string;
  /**
   * The path (subordinate to the group) where the workflow will exist.
   * (optionally) use "" to override the landing page, will need to also update the route path
   */
  path: string;
  /**
   * The routes correspond to different tasks in a workflow (i.e. for AWS EC2,
   * users could `reboot instance` or `terminate instance`).
   */
  routes: unknown;
}

interface WorkflowShortlinkConfiguration {
  /**
   * (Optional) property to enable a workflow to utilize short linking of its state.
   * This property is required to show the short link generator in the header
   */
  shortLink?: boolean;
}

export interface Workflow extends BaseWorkflowConfiguration, WorkflowShortlinkConfiguration {
  /**
   * An optional property that is set via the config and allows for the display of an icon given a path,
   * this will override the default avatar.
   * { path: string }
   */
  icon?: WorkflowIcon;
  /**
   * Configured routes allow for the optional properties of `trending` (whether to display
   * on homepage) and `componentProps` which allow the passing of workflow/route
   * specific props.
   */
  routes: ConfiguredRoute[];
}

export interface WorkflowConfiguration
  extends BaseWorkflowConfiguration,
    WorkflowShortlinkConfiguration {
  shortLink?: boolean;
  routes: {
    [key: string]: Route;
  };
}

export interface Route {
  component: React.FC<any>;
  description: string;
  displayName?: string;
  /** (optionally) use "" to override the landing page */
  path: string;
  /** Properties required by the Component that are set only via the config. */
  requiredConfigProps?: string[];
  /** Is the workflow discoverable via search and drawer navigation. This defaults to false. */
  hideNav?: boolean;
  /**
   * The feature flag used to determine if the route should be registered.
   *
   * If this is not set the route will always be registered.
   */
  featureFlag?: string;

  layoutProps?: LayoutProps;
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
      <>Failed to load {workflow.displayName} workflow. Please contact the developer</>
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
        <Grid container direction="column" justifyContent="center" alignItems="center">
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
