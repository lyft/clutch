import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";

import HelloWorld from "./hello-world";

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
        layoutProps: {
          {{- if .IsWizardTemplate}}
          variant: "wizard",
          {{- else}}
          variant: "standard",
          {{- end}}
        },
      },
    },
  };
};

export default register;
