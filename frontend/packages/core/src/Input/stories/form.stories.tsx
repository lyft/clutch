import * as React from "react";
import { action } from "@storybook/addon-actions";
import type { Meta } from "@storybook/react";

import { Button } from "../../button";
import { Form, FormRow } from "../form";
import { TextField } from "../text-field";

export default {
  title: "Core/Input/Form",
  component: Form,
} as Meta;

const Template = ({ hasRow = false }: { hasRow: boolean }) => {
  const children = (
    <>
      <TextField />
      <TextField />
    </>
  );
  return (
    <Form
      onSubmit={e => {
        e.preventDefault();
        action("onSubmit event")(e);
      }}
    >
      {hasRow && <FormRow>{children}</FormRow>}
      {children}
      <Button text="Submit Form" type="submit" />
    </Form>
  );
};

export const Primary = Template.bind({});

export const WithRows = Template.bind({});
WithRows.args = {
  hasRow: true,
};
