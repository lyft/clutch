import React from "react";
import { MemoryRouter } from "react-router";
import type { Meta } from "@storybook/react";

import { ApplicationContext } from "../../Contexts/app-context";
import Drawer from "../drawer";
import { Primary as HeaderStory } from "./header.stories";

export default {
  title: "Core/AppLayout/Drawer",
  component: Drawer,
  decorators: [
    () => (
      <MemoryRouter>
        <Drawer />
      </MemoryRouter>
    ),
    StoryFn => {
      return (
        <ApplicationContext.Provider value={{
          workflows: [
            {
              developer: { name: "Lyft", contactUrl: "mailto:hello@clutch.sh" },
              displayName: "EC2",
              group: "AWS",
              path: "ec2",
              routes: [
                {
                  component: () => <></>,
                  componentProps: { resolverType: "clutch.aws.ec2.v1.Instance" },
                  description: "Terminate an EC2 instance.",
                  displayName: "Terminate Instance",
                  path: "instance/terminate",
                  requiredConfigProps: ["resolverType"],
                  trending: true,
                },
                {
                  component: () => <></>,
                  componentProps: { resolverType: "clutch.aws.ec2.v1.AutoscalingGroup" },
                  description: "Resize an autoscaling group.",
                  displayName: "Resize Autoscaling Group",
                  path: "asg/resize",
                  requiredConfigProps: ["resolverType"],
                },
              ],
            },
          ],
        }}>
          <StoryFn />
        </ApplicationContext.Provider>
      );
    },
  ],
} as Meta;

export const Primary = () => <Drawer />;
