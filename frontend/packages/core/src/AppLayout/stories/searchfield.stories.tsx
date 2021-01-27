import * as React from "react";
import { MemoryRouter } from "react-router";
import styled from "@emotion/styled";
import { Box, Grid as MuiGrid } from "@material-ui/core";
import type { Meta } from "@storybook/react";

import { ApplicationContext } from "../../Contexts/app-context";
import SearchField from "../search";

export default {
  title: "Core/AppLayout/Search Field",
  component: SearchField,
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
          value={{
            workflows: [
              {
                developer: { name: "Lyft", contactUrl: "mailto:hello@clutch.sh" },
                displayName: "EC2",
                group: "AWS",
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

const Grid = styled(MuiGrid)({
  height: "64px",
  backgroundColor: "#131C5F",
});

const Template = () => (
  <Grid container alignItems="center" justify="center">
    <Box>
      <SearchField />
    </Box>
  </Grid>
);

export const Primary = Template.bind({});
