import * as React from "react";
import { Grid } from "@material-ui/core";
import type { Meta } from "@storybook/react";

import { UserInformation } from "../user";

export default {
  title: "Core/AppLayout/User Information",
  component: UserInformation,
  parameters: {
    backgrounds: {
      default: "header blue",
      values: [{ name: "header blue", value: "#131C5F" }],
    },
  },
} as Meta;

const Template = () => (
  <Grid container alignItems="center" justify="center">
    <UserInformation />
  </Grid>
);
export const Primary = Template.bind({});
