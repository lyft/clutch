import * as React from "react";
import type { Meta } from "@storybook/react";

import { Grid } from "../../Layout";
import { styled } from "../../Utils";
import type { NotificationsProp } from "../notifications";
import Notifications from "../notifications";

export default {
  title: "Core/AppLayout/Notifications",
  component: Notifications,
} as Meta;

const StyledGrid = styled(Grid)({
  height: "64px",
  backgroundColor: "#131C5F",
});

const Template = (props: NotificationsProp) => (
  <StyledGrid container alignItems="center" justifyContent="center">
    <Notifications {...props} />
  </StyledGrid>
);

export const Primary = Template.bind({});
Primary.args = {
  data: [{ value: "New K8s workflow!" }, { value: "Clutch v1.18 release" }],
};
