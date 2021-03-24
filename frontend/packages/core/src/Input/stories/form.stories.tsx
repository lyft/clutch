import * as React from "react";
import { action } from "@storybook/addon-actions";
import type { Meta } from "@storybook/react";

import { Button } from "../../button";
import type { FormProps } from "../form";
import Form from "../form";
import { TextField } from "../text-field";

export default {
  title: "Core/Input/Form",
  component: Form,
} as Meta;

const Template = (props: FormProps) => (
  <Form
    onSubmit={e => {
      e.preventDefault();
      action("onSubmit event")(e);
    }}
    {...props}
  >
    <TextField />
    <TextField />
    <Button text="Submit Form" type="submit" />
  </Form>
);

export const Primary = Template.bind({});

export const Horizontal = Template.bind({});
Horizontal.args = {
  direction: "row",
};
