import * as React from "react";
import type { Meta } from "@storybook/react";

import { Grid } from "../../Layout";
import type { NotificationsProp } from "../notifications";
import NotificationsComponent from "../notifications";

export default {
  title: "Core/AppLayout/Notifications",
  component: NotificationsComponent,
} as Meta;

const Template = (props: NotificationsProp) => (
  <Grid container alignItems="center" justifyContent="center">
    <NotificationsComponent {...props} />
  </Grid>
);

export const Notifications = Template.bind({});
Notifications.args = {
  data: [{ value: "New K8s workflow!" }, { value: "Clutch v1.18 release" }],
};
