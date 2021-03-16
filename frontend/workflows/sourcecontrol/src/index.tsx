import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";
import type { WizardChild } from "@clutch-sh/wizard";

import CreateRepository from "./create-repository";

interface RepositoryConfigProps {
  options: {
    [option: string]: boolean;
  };
}

export interface WorkflowProps extends BaseWorkflowProps, RepositoryConfigProps {}
export interface RepostioryChild extends WizardChild, RepositoryConfigProps {}

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "scm",
    group: "Source Control",
    displayName: "Source Control",
    routes: {
      createRepository: {
        path: "createRepository",
        component: CreateRepository,
        displayName: "Create Repository",
        description: "Create a new repository.",
      },
    },
  };
};

export default register;
