import * as React from "react";
import { MemoryRouter } from "react-router";
import styled from "@emotion/styled";
import { Box as MuiBox } from "@material-ui/core";
import type { Meta } from "@storybook/react";

import { ApplicationContext } from "../../Contexts/app-context";
import Drawer from "../drawer";

export default {
  title: "Core/AppLayout/Drawer",
  component: Drawer,
  decorators: [
    Sidebar => (
      <MemoryRouter>
        <Sidebar />
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
          }}
        >
          <StoryFn />
        </ApplicationContext.Provider>
      );
    },
  ],
} as Meta;

const Box = styled(MuiBox)({
  ".MuiDrawer-root div[class*='MuiDrawer-paper']": {
    top: 0,
  },
});

const Template = () => (
  <Box>
    <Drawer />
  </Box>
);

export const Primary = Template.bind({});
