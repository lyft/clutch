import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";

import HelloWorld from "./{{ .HelloWorldModule }}";

export interface WorkflowProps extends BaseWorkflowProps {}

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "{{ .DeveloperName }}",
      contactUrl: "mailto:{{ .DeveloperEmail }}",
    },
    path: "{{ .URLRoot }}",
    group: "{{ .Name }}",
    displayName: "{{ .Name }}",
    routes: {
      landing: {
        path: "{{ .URLPath }}",
        displayName: "{{ .Name }}",
        description: "{{ .Description }}.",
        component: HelloWorld,
      },
    },
  };
};

export default register;
