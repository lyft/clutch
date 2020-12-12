import * as React from "react";
import styled from "@emotion/styled";
import { Grid as MuiGrid } from "@material-ui/core";
import type { Meta } from "@storybook/react";

import type { UserNotficationsProp } from "../notifications";
import { UserNotifications } from "../notifications";

export default {
  title: "Core/AppLayout/User Notifications",
  component: UserNotifications,
} as Meta;

const Grid = styled(MuiGrid)({
  height: "64px",
  backgroundColor: "#131C5F",
});

const Template = (props: UserNotficationsProp) => (
  <Grid container alignItems="center" justify="center">
    <UserNotifications {...props} />
  </Grid>
);

export const Primary = Template.bind({});
Primary.args = {
  data: [{ value: "New K8s workflow!" }, { value: "Clutch v1.18 release" }],
};
