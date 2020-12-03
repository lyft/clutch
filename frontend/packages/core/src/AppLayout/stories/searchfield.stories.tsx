import * as React from "react";
import { MemoryRouter } from "react-router";
import type { Meta } from "@storybook/react";

import ResizeAutoscalingGroup from "../../../../../workflows/ec2/src/resize-asg";
import TerminateInstance from "../../../../../workflows/ec2/src/terminate-instance";
import { ApplicationContext } from "../../Contexts/app-context";
import SearchField from "../search";

export default {
  title: "Core/AppLayout/Search Field",
  component: SearchField,
  decorators: [
    () => (
      <MemoryRouter>
        <SearchField />
      </MemoryRouter>
    ),
    StoryFn => {
      return (
        <ApplicationContext.Provider
          value={{
            workflows: [
              {
                developer: { name: "Lyft", contactUrl: "mailto:hello@clutch.sh" },
                displayName: "EC2",
                group: "AWS",
                path: "ec2",
                routes: [
                  {
                    component: TerminateInstance,
                    componentProps: { resolverType: "clutch.aws.ec2.v1.Instance" },
                    description: "Terminate an EC2 instance.",
                    displayName: "Terminate Instance",
                    path: "instance/terminate",
                    requiredConfigProps: ["resolverType"],
                    trending: true,
                  },
                  {
                    component: ResizeAutoscalingGroup,
                    componentProps: { resolverType: "clutch.aws.ec2.v1.AutoscalingGroup" },
                    description: "Resize an autoscaling group.",
                    displayName: "Resize Autoscaling Group",
                    path: "asg/resize",
                    requiredConfigProps: ["resolverType"],
                  },
                ],
              },
            ],
          }}
        >
          <StoryFn />
        </ApplicationContext.Provider>
      );
    },
  ],
  parameters: {
    backgrounds: {
      default: "header blue",
      values: [{ name: "header blue", value: "#131C5F" }],
    },
  },
} as Meta;

export const Primary: React.FC<{}> = () => <SearchField />;
