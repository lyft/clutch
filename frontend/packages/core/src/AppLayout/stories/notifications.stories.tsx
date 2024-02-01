import * as React from "react";
import { Grid as MuiGrid, Theme } from "@mui/material";
import type { Meta } from "@storybook/react";

import styled from "../../styled";
import type { NotificationsProp } from "../notifications";
import NotificationsComponent from "../notifications";

export default {
  title: "Core/AppLayout/Notifications",
  component: NotificationsComponent,
} as Meta;

const Grid = styled(MuiGrid)(({ theme }: { theme: Theme }) => ({
  height: "64px",
  backgroundColor: theme.palette.primary[900],
}));

const Template = (props: NotificationsProp) => (
  <Grid container alignItems="center" justifyContent="center">
    <NotificationsComponent {...props} />
  </Grid>
);

export const Notifications = Template.bind({});
Notifications.args = {
  data: [{ value: "New K8s workflow!" }, { value: "Clutch v1.18 release" }],
};
