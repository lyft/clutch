import * as React from "react";
import { Grid as MuiGrid } from "@mui/material";
import type { Meta } from "@storybook/react";

import styled from "../../styled";
import type { UserInformationProps } from "../user";
import { UserInformation as UserInformationComponent } from "../user";

export default {
  title: "Core/AppLayout/User Information",
  component: UserInformationComponent,
} as Meta;

const Grid = styled(MuiGrid)(({ theme }) => ({
  height: "64px",
  backgroundColor: theme.palette.primary[900],
}));

const Template = (props: UserInformationProps) => (
  <Grid container alignItems="center" justifyContent="center">
    <UserInformationComponent {...props} />
  </Grid>
);
export const UserInformation = Template.bind({});
UserInformation.args = {
  data: [{ value: "Dashboard" }, { value: "Settings" }],
  user: "fooBar@example.com",
};
