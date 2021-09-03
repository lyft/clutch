import React from "react";
import ReactDOM from "react-dom";
import { ClutchApp } from "@clutch-sh/core";

import availableWorkflows from "./workflows";

import "./index.css";

const config = require("./clutch.config.js");

ReactDOM.render(
  <ClutchApp workflows={availableWorkflows} gatewayConfig={config} />,
  document.getElementById("root")
);
