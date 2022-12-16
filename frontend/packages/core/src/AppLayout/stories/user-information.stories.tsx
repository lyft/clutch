import * as React from "react";
import type { Meta } from "@storybook/react";

import { Grid } from "../../Layout";
import { styled } from "../../Utils";
import type { UserInformationProps } from "../user";
import { UserInformation } from "../user";

export default {
  title: "Core/AppLayout/User Information",
  component: UserInformation,
} as Meta;

const StyledGrid = styled(Grid)({
  height: "64px",
  backgroundColor: "#131C5F",
});

const Template = (props: UserInformationProps) => (
  <StyledGrid container alignItems="center" justifyContent="center">
    <UserInformation {...props} />
  </StyledGrid>
);
export const Primary = Template.bind({});
Primary.args = {
  data: [{ value: "Dashboard" }, { value: "Settings" }],
  user: "fooBar@example.com",
};
