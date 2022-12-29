import * as React from "react";
import type { Meta } from "@storybook/react";

import { SEVERITIES } from "../../Assets/global";
import { Alert as AlertComponent, AlertProps } from "../alert";

export default {
  title: "Core/Feedback/Alert",
  component: AlertComponent,
  argTypes: {
    title: {
      control: {
        type: "text",
      },
    },
    severity: {
      options: SEVERITIES,
      control: {
        type: "select",
      },
    },
  },
} as Meta;

const Template = (props: AlertProps) => <AlertComponent {...props}>This is a note</AlertComponent>;

export const Alert = Template.bind({});
Alert.args = {
  title: "A Title",
  severity: "success",
};
