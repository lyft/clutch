import React from "react";
import ReactDOM from "react-dom";
import { ClutchApp } from "@clutch-sh/core";

import registeredWorkflows from "./workflows";

import "./index.css";

import config from "./clutch.config";

ReactDOM.render(
  <ClutchApp availableWorkflows={registeredWorkflows} configuration={config} />,
  document.getElementById("root")
);
