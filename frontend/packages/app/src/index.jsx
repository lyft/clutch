import React from "react";
import ReactDOM from "react-dom";
import { ClutchApp } from "@clutch-sh/core";

import registeredWorkflows from "./workflows";

import "./index.css";

const config = require("./clutch.config");

const banners = {
  header: {
    title: "this is a header banner",
    message: "this is message header",
    dismissed: false,
  },
  perWorkflow: {
    Envoy: {
      title: "this is a per workout banner fro envoy",
      message: "this is message banner fro envoy",
      dismissed: false,
    },
    K8s: {
      title: "this is a per workout banner for k8s",
      message: "this is message banner for k8s",
      dismissed: false,
    },
  },
  multiWorkflow: {
    title: "multi title",
    message: "message for multi workflow",
    dismissed: false,
    workflows: ["EC2", "Envoy"],
  },
};

ReactDOM.render(
  <ClutchApp
    availableWorkflows={registeredWorkflows}
    configuration={config}
    appConfiguration={{ banners }}
  />,
  document.getElementById("root")
);
