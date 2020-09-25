import type { BaseWorkflowProps, NoteConfig, WorkflowConfiguration } from "@clutch-sh/core";
import type { WizardChild } from "@clutch-sh/wizard";

import DeletePod from "./delete-pod";
import ResizeHPA from "./resize-hpa";
import DescribeService from "./describe-service";
import ListServices from "./list-services";

interface ResolverConfigProps {
  resolverType: string;
}

interface ConfirmConfigProps {
  notes?: NoteConfig[];
}

export interface WorkflowProps extends BaseWorkflowProps, ResolverConfigProps, ConfirmConfigProps {}
export interface ResolverChild extends WizardChild, ResolverConfigProps {}
export interface ConfirmChild extends WizardChild, ConfirmConfigProps {}

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "k8s",
    group: "K8s",
    displayName: "K8s",
    routes: {
      deletePod: {
        path: "pod/delete",
        displayName: "Delete Pod",
        description: "Delete a K8s pod.",
        component: DeletePod,
        requiredConfigProps: ["resolverType"],
      },
      resizeHPA: {
        path: "hpa/resize",
        displayName: "Resize HPA",
        description: "Resize a horizontal autoscaler.",
        component: ResizeHPA,
        requiredConfigProps: ["resolverType"],
      },
      describeService: {
        path: "svc/describe",
        displayName: "Describe Service",
        description: "Describe a running K8s service.",
        component: DescribeService,
        requiredConfigProps: ["resolverType"],
      },
      listServices: {
        path: "svc/list",
        displayName: "List Services",
        description: "View a list of running K8s services.",
        component: ListServices,
        requiredConfigProps: ["resolverType"],
      },
    },
  };
};

export default register;
