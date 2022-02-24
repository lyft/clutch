import React from "react";
import { action } from "@storybook/addon-actions";
import type { Meta } from "@storybook/react";

import type { ToastProps } from "../toast";
import Toast from "../toast";

export default {
  title: "Core/Feedback/Toast",
  component: Toast,
  argTypes: {
    severity: {
      defaultValue: "warning",
      description: "Severity Level for Component",
      options: ["success", "info", "warning", "error"],
      control: { type: "select" },
      table: {
        type: { summary: "enum" },
        defaultValue: { summary: "warning" },
      },
    },
    onClose: {
      description: "On Close handler for component",
      options: ["enabled", "disabled"],
      defaultValue: "disabled",
      control: { type: "radio" },
      mapping: {
        enabled: action("on-close"),
        disabled: undefined,
      },
      table: {
        type: { summary: "function" },
      },
    },
    autoHideDuration: {
      description: "Auto-hide duration for Toast",
      defaultValue: null,
      disable: true,
      table: {
        type: { summary: "number" },
        defaultValue: { summary: 6000 },
      },
    },
    title: {
      description: "Title of Toast",
      table: {
        type: { summary: "React.ReactNode" },
      },
    },
    anchorOrigin: {
      control: "object",
      description: "Location of Toast",
      defaultValue: { vertical: "bottom", horizontal: "right" },
      table: {
        type: { summary: "object" },
        defaultValue: { summary: `{ vertical: "bottom", horizontal: "right" }` },
      },
    },
  },
} as Meta;

const Template = ({ severity, ...props }: ToastProps) => (
  <Toast autoHideDuration={null} severity={severity} {...props}>
    Informational {severity} Toast
  </Toast>
);

export const Primary = Template.bind({});
Primary.args = {};

export const WithTitle = Template.bind({});
WithTitle.args = {
  ...Primary.args,
  title: "Title Message",
};
