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
    link: "localhost:3000",
  },
  perWorkflow: {
    Envoy: {
      title: "this is a per workout banner fro envoy",
      message: "this is message banner fro envoy",
      link: "localhost:3000",
    },
    K8s: {
      title: "this is a per workout banner for k8s",
      message: "this is message banner for k8s",
      link: "localhost:3000",
    },
    test: {
      title: "this is a per workout banner for test",
      message: "this is message banner for test",
      link: "localhost:3000",
    },
  },
  multiWorkflow: {
    title: "multi title",
    message: "message for multi workflow",
    workflows: ["EC2", "Envoy"],
    link: "localhost:3000",
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
