import * as React from "react";
import { action } from "@storybook/addon-actions";
import type { Meta } from "@storybook/react";

import { Button } from "../../button";
import Form from "../form";
import { TextField } from "../text-field";

export default {
  title: "Core/Input/Form",
  component: Form,
} as Meta;

const Template = () => (
  <Form
    onSubmit={e => {
      e.preventDefault();
      action("onSubmit event")(e);
    }}
  >
    <TextField />
    <Button text="Submit Form" type="submit" />
  </Form>
);

export const Primary = Template.bind({});
Primary.args = {
  selected: false,
};
