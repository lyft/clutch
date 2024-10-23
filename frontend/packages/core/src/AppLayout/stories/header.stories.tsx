import * as React from "react";
import { MemoryRouter } from "react-router-dom";
import type { Meta, StoryObj } from "@storybook/react";

import { ApplicationContext, UserPreferencesProvider } from "../../Contexts";
import { Header } from "../header";

const workflows = [
  {
    developer: { name: "Lyft", contactUrl: "mailto:hello@clutch.sh" },
    displayName: "EC2",
    group: "AWS",
    icon: { path: "" },
    path: "ec2",
    routes: [
      {
        component: () => <div>Terminate Instance</div>,
        componentProps: { resolverType: "clutch.aws.ec2.v1.Instance" },
        description: "Terminate an EC2 instance.",
        displayName: "Terminate Instance",
        path: "instance/terminate",
        requiredConfigProps: ["resolverType"],
        trending: true,
      },
      {
        component: () => <div>Resize ASG</div>,
        componentProps: { resolverType: "clutch.aws.ec2.v1.AutoscalingGroup" },
        description: "Resize an autoscaling group.",
        displayName: "Resize Autoscaling Group",
        path: "asg/resize",
        requiredConfigProps: ["resolverType"],
      },
    ],
  },
];

const meta: Meta<typeof Header> = {
  title: "Core/Layout/Header",
  component: Header,
  argTypes: {
    logo: {
      control: "text",
    },
  },
  decorators: [
    StoryFn => (
      <ApplicationContext.Provider
        // eslint-disable-next-line react/jsx-no-constructed-context-values
        value={{
          workflows,
        }}
      >
        <UserPreferencesProvider>
          <StoryFn />
        </UserPreferencesProvider>
      </ApplicationContext.Provider>
    ),
  ],
};

export default meta;
type Story = StoryObj<typeof Header>;

export const Basic: Story = {
  render: args => (
    <MemoryRouter>
      <Header {...args} />
    </MemoryRouter>
  ),
};
