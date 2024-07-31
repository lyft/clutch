import * as React from "react";
import { MemoryRouter } from "react-router";
import styled from "@emotion/styled";
import { Box, Grid as MuiGrid, Theme } from "@mui/material";
import type { Meta } from "@storybook/react";

import { ApplicationContext } from "../../Contexts/app-context";
import { THEME_VARIANTS } from "../../Theme/colors";
import SearchFieldComponent from "../search";

export default {
  title: "Core/AppLayout/Search Field",
  component: SearchFieldComponent,
  decorators: [
    Search => {
      return (
        <MemoryRouter>
          <Search />
        </MemoryRouter>
      );
    },
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
            ],
          }}
        >
          <StoryFn />
        </ApplicationContext.Provider>
      );
    },
  ],
} as Meta;

const Grid = styled(MuiGrid)(({ theme }: { theme: Theme }) => ({
  height: "64px",
  backgroundColor:
    theme.palette.mode === THEME_VARIANTS.light
      ? theme.palette.primary[900]
      : theme.palette.headerGradient,
}));

const Template = () => (
  <Grid container alignItems="center" justifyContent="center">
    <Box>
      <SearchFieldComponent />
    </Box>
  </Grid>
);

export const SearchField = Template.bind({});
