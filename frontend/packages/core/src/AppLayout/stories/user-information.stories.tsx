import * as React from "react";
import styled from "@emotion/styled";
import { Grid as MuiGrid } from "@material-ui/core";
import type { Meta } from "@storybook/react";

import type { UserInformationProps } from "../user";
import { UserInformation } from "../user";

export default {
  title: "Core/AppLayout/User Information",
  component: UserInformation,
} as Meta;

const Grid = styled(MuiGrid)({
  height: "64px",
  backgroundColor: "#131C5F",
});

const Template = (props: UserInformationProps) => (
  <Grid container alignItems="center" justify="center">
    <UserInformation {...props} />
  </Grid>
);
export const Primary = Template.bind({});
Primary.args = {
  data: [{ value: "Dashboard" }, { value: "Settings" }],
  user: "fooBar@gmail.com",
};
