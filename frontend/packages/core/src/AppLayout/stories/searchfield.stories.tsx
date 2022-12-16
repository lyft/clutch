import * as React from "react";
import { MemoryRouter } from "react-router";
import { Box } from "@mui/material";
import type { Meta } from "@storybook/react";

import { ApplicationContext } from "../../Contexts/app-context";
import { Grid } from "../../Layout";
import { styled } from "../../Utils";
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

const StyledGrid = styled(Grid)({
  height: "64px",
  backgroundColor: "#131C5F",
});

const Template = () => (
  <StyledGrid container alignItems="center" justifyContent="center">
    <Box>
      <SearchField />
    </Box>
  </StyledGrid>
);

export const Primary = Template.bind({});
