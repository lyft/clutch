import * as React from "react";
import type { Meta } from "@storybook/react";

import type { MultiSelectProps } from "../select";
import { MultiSelect } from "../select";

export default {
  title: "Core/Input/MultiSelect",
  component: MultiSelect,
  argTypes: {
    onChange: { action: "onChange event" },
    selectOptions: { control: { type: "object" } },
  },
} as Meta;

const Template = (props: MultiSelectProps) => <MultiSelect name="storybookDemo" {...props} />;

export const Default = Template.bind({});
Default.args = {
  label: "My Label",
  selectOptions: [
    {
      label: "Option 1",
    },
    {
      label: "Option 2",
      value: "Other Value",
    },
  ],
};

export const WithStartAdornment = Template.bind({});
WithStartAdornment.args = {
  ...Default.args,
  selectOptions: [
    {
      label: "Option 1",
      startAdornment: <img src="https://clutch.sh/img/microsite/logo.svg" alt="logo" />,
    },
  ],
};

export const WithGrouping = Template.bind({});
WithGrouping.args = {
  ...Default.args,
  selectOptions: [
    {
      label: "Option 1",
    },
    {
      label: "Group 1",
      group: [
        {
          label: "Sub Option 1",
        },
        {
          label: "Sub Option 2",
        },
      ],
    },
    {
      label: "Group 2",
      group: [
        {
          label: "Sub Option 3",
        },
        {
          label: "Sub Option 4",
        },
      ],
    },
    {
      label: "Option 2",
    },
  ],
};
