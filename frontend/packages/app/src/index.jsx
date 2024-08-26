import React from "react";
import ReactDOM from "react-dom";
import { ClutchApp } from "@clutch-sh/core";

import registeredWorkflows from "./workflows";

import "./index.css";

const config = require("./clutch.config");

const banners = {
  header: {
    message: "this is message header",
    linkText: "infra docs",
    link: "https://infradocs.lyft.net/index.html",
    severity: "error",
  },
  perWorkflow: {
    Envoy: {
      title: "this is a per workout banner fro envoy",
      message: "this is message banner fro envoy",
      linkText: "infra docs",
      link: "https://infradocs.lyft.net/index.html",
      severity: "success",
    },
    K8s: {
      title: "this is a per workout banner for k8s",
      message: "this is message banner for k8s",
      linkText: "infra docs",
      link: "https://infradocs.lyft.net/index.html",
      severity: "warning",
    },
  },
  multiWorkflow: {
    title: "multi title",
    message: "message for multi workflow",
    workflows: ["EC2", "Envoy"],
    severity: "info",
    linkText: "infra docs",
    link: "https://infradocs.lyft.net/index.html",
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
