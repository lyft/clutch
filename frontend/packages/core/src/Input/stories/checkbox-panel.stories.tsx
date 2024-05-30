import * as React from "react";
import type { Meta } from "@storybook/react";

import type { CheckboxPanelProps } from "../checkbox";
import { CheckboxPanel as CheckboxPanelComponent } from "../checkbox";

export default {
  title: "Core/Input/CheckboxPanel",
  component: CheckboxPanelComponent,
  argTypes: {
    onChange: { action: "onChange event" },
  },
} as Meta;

const Template = (props: CheckboxPanelProps) => <CheckboxPanelComponent {...props} />;

export const CheckboxPanel = Template.bind({});
CheckboxPanel.args = {
  header: "Select all that apply:",
  options: {
    "Option 1": false,
    "Option 2": false,
    "Option 3": false,
  },
};

export const WithClearOption = Template.bind({});
WithClearOption.args = {
  header: "Select all that apply:",
  options: {
    "Option 1": false,
    "Option 2": false,
    "Option 3": false,
  },
  onClear: () => null,
};
