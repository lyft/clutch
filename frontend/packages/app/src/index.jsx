import React from "react";
import ReactDOM from "react-dom";
import { ClutchApp } from "@clutch-sh/core";

import registeredWorkflows from "./workflows";

import "./index.css";

const config = require("./clutch.config");

ReactDOM.render(
  <ClutchApp availableWorkflows={registeredWorkflows} configuration={config} />,
  document.getElementById("root")
);
