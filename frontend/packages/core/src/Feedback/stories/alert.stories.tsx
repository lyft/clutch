import * as React from "react";
import type { Meta } from "@storybook/react";

import type { AlertProps } from "../alert";
import { Alert } from "../alert";

export default {
  title: "Core/Feedback/Alert",
  component: Alert,
} as Meta;

const Template = (props: AlertProps) => <Alert {...props}>This is a note</Alert>;

export const WithTitle = Template.bind({});
WithTitle.args = {
  title: "A Title",
};

export const Success = Template.bind({});
Success.args = {
  severity: "success",
  title: "A title",
};

export const Error = Template.bind({});
Error.args = {
  severity: "error",
};

export const Info = Template.bind({});
Info.args = {
  severity: "info",
};

export const Warning = Template.bind({});
Warning.args = {
  severity: "warning",
};
