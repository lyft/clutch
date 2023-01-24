import * as React from "react";
import SearchIcon from "@mui/icons-material/Search";
import type { Meta } from "@storybook/react";

import type { TextFieldProps } from "../text-field";
import { TextField } from "../text-field";

export default {
  title: "Core/Input/TextField",
  component: TextField,
  argTypes: {
    color: {
      options: ["primary", "secondary", "error", "info", "success", "warning"],
      control: { type: "select" },
      defaultValue: "primary",
    },
  },
} as Meta;

const Template = (props: TextFieldProps) => <TextField {...props} />;

export const Primary = Template.bind({});
Primary.args = {
  label: "My Label",
  color: "primary",
  error: false,
  placeholder: "This is a placeholder, start typing",
};

export const Disabled = Template.bind({});
Disabled.args = {
  ...Primary.args,
  disabled: true,
};

export const Error = Template.bind({});
Error.args = {
  ...Primary.args,
  error: true,
  helperText: "There was a problem!",
};

export const WithoutLabel = Template.bind({});
WithoutLabel.args = {
  ...Primary.args,
  label: null,
};

export const MultipleLines = Template.bind({});
MultipleLines.args = {
  ...Primary.args,
  multiline: true,
  defaultValue: "This is\nan example\nof multiline content",
};

export const WithEndAdornment = Template.bind({});
WithEndAdornment.args = {
  ...Primary.args,
  defaultValue: "Search",
  endAdornment: <SearchIcon />,
};

const autoComplete = value => {
  return new Promise((resolve, reject) => {
    resolve({
      results: [
        { id: "clutch", label: "" },
        { id: "clutch-auto", label: "" },
        { id: "clutch-autocomplete", label: "" },
      ],
    });
    reject(new Error("Something bad happened"));
  });
};

export const Autocomplete = Template.bind({});
Autocomplete.args = {
  ...Primary.args,
  placeholder: "Search for `clutch`",
  autocompleteCallback: autoComplete,
};
