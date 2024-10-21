import React from "react";
import { MemoryRouter } from "react-router";
import type { Meta } from "@storybook/react";

import { Header } from "../../AppLayout/header";
import { ApplicationContext, UserPreferencesProvider } from "../../Contexts";
import { NPSHeader } from "..";

export default {
  title: "Core/NPS/Header",
  component: NPSHeader,
} as Meta;

const Template = () => (
  <ApplicationContext.Provider
    // eslint-disable-next-line react/jsx-no-constructed-context-values
    value={{
      workflows: [
        {
          developer: { name: "Lyft", contactUrl: "mailto:hello@clutch.sh" },
          displayName: "EC2",
          group: "AWS",
          path: "ec2",
          routes: [
            {
              // eslint-disable-next-line react/no-unstable-nested-components
              component: () => <div>Terminate Instance</div>,
              componentProps: { resolverType: "clutch.aws.ec2.v1.Instance" },
              description: "Terminate an EC2 instance.",
              displayName: "Terminate Instance",
              path: "instance/terminate",
              requiredConfigProps: ["resolverType"],
              trending: true,
            },
            {
              // eslint-disable-next-line react/no-unstable-nested-components
              component: () => <div>Resize ASG</div>,
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
    <UserPreferencesProvider>
      <MemoryRouter>
        <Header enableNPS />
      </MemoryRouter>
    </UserPreferencesProvider>
  </ApplicationContext.Provider>
);

export const Primary = Template.bind({});
