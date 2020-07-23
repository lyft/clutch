import React from "react";
import { ClutchApp } from "@clutch-sh/core";
import ReactDOM from "react-dom";

import registeredWorkflows from "./workflows";
import "./index.css";

const config = require("./clutch.config.js");

ReactDOM.render(
    <ClutchApp availableWorkflows={registeredWorkflows} configuration={config} />,
    document.getElementById("root")
);
