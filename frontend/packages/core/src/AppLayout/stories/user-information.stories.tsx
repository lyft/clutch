import * as React from "react";
import type { Meta } from "@storybook/react";

import { Grid } from "../../Layout";
import type { UserInformationProps } from "../user";
import { UserInformation as UserInformationComponent } from "../user";

export default {
  title: "Core/AppLayout/User Information",
  component: UserInformationComponent,
} as Meta;

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
