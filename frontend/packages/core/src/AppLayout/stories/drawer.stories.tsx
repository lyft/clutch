import * as React from "react";
import { MemoryRouter } from "react-router";
import { Grid } from "@mui/material";
import type { Meta } from "@storybook/react";

import { UserPreferencesProvider } from "../../Contexts";
import { ApplicationContext } from "../../Contexts/app-context";
import Drawer from "../drawer";
import Header from "../header";

export default {
  title: "Core/Layout/Drawer",
  component: Drawer,
  decorators: [
    StoryFn => (
      <MemoryRouter>
        <StoryFn />
      </MemoryRouter>
    ),
    StoryFn => {
      return (
        <ApplicationContext.Provider
          // eslint-disable-next-line react/jsx-no-constructed-context-values
          value={{
            workflows: [
              {
                developer: { name: "Lyft", contactUrl: "mailto:hello@clutch.sh" },
                displayName: "EC2",
                group: "AWS",
                path: "ec2",
                icon: { path: "" },
                routes: [
                  {
                    component: () => null,
                    componentProps: { resolverType: "clutch.aws.ec2.v1.Instance" },
                    description: "Terminate an EC2 instance.",
                    displayName: "Terminate Instance",
                    path: "instance/terminate",
                    requiredConfigProps: ["resolverType"],
                    trending: true,
                  },
                  {
                    component: () => null,
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
    layout: "fullscreen",
  },
} as Meta;

export const Primary = () => <Drawer />;

export const WithHeader = () => (
  <UserPreferencesProvider>
    <Grid container direction="column">
      <Header />
      <Drawer />
    </Grid>
  </UserPreferencesProvider>
);
