import React from "react";
import ReactDOM from "react-dom";
import Bugsnag from "@bugsnag/js";
import BugsnagPluginReact from "@bugsnag/plugin-react";
import { ClutchApp } from "@clutch-sh/core";

import registeredWorkflows from "./workflows";

import "./index.css";

const config = require("./clutch.config");

let root = <ClutchApp availableWorkflows={registeredWorkflows} configuration={config} />;

if (process.env.REACT_APP_CREDENTIALS_BUGSNAG_API_TOKEN) {
  Bugsnag.start({
    apiKey: process.env.REACT_APP_CREDENTIALS_BUGSNAG_API_TOKEN,
    plugins: [new BugsnagPluginReact()],
    releaseStage: process.env.APPLICATION_ENV,
  });
  const ErrorBoundary = Bugsnag.getPlugin("react").createErrorBoundary(React);
  root = <ErrorBoundary>{root}</ErrorBoundary>;
}

ReactDOM.render(root, document.getElementById("root"));
