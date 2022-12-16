import * as React from "react";
import { Grid as MuiGrid } from "@mui/material";
import type { Meta } from "@storybook/react";

import { styled } from "../../Utils";
import type { NotificationsProp } from "../notifications";
import Notifications from "../notifications";

export default {
  title: "Core/AppLayout/Notifications",
  component: Notifications,
} as Meta;

const Grid = styled(MuiGrid)({
  height: "64px",
  backgroundColor: "#131C5F",
});

const Template = (props: NotificationsProp) => (
  <Grid container alignItems="center" justifyContent="center">
    <Notifications {...props} />
  </Grid>
);

export const Primary = Template.bind({});
Primary.args = {
  data: [{ value: "New K8s workflow!" }, { value: "Clutch v1.18 release" }],
};
