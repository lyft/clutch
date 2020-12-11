import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";
import type { WizardChild } from "@clutch-sh/wizard";

import RebootInstance from "./reboot-instance";
import ResizeAutoscalingGroup from "./resize-asg";
import TerminateInstance from "./terminate-instance";

interface ResolverConfigProps {
  resolverType: string;
}

interface ConfirmConfigProps {
  note?: string;
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
    path: "ec2",
    group: "AWS",
    displayName: "EC2",
    routes: {
      terminateInstance: {
        path: "instance/terminate",
        displayName: "Terminate Instance",
        description: "Terminate an EC2 instance.",
        component: TerminateInstance,
        requiredConfigProps: ["resolverType"],
      },
      rebootInstance: {
        path: "instance/reboot",
        displayName: "Reboot Instance",
        description: "Reboot an EC2 Instance",
        component: RebootInstance,
        requiredConfigProps: ["resolverType"],
      },
      resizeAutoscalingGroup: {
        path: "asg/resize",
        displayName: "Resize Autoscaling Group",
        description: "Resize an autoscaling group.",
        component: ResizeAutoscalingGroup,
        requiredConfigProps: ["resolverType"],
      },
    },
  };
};

export default register;
