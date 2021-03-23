import * as React from "react";
import type { Meta } from "@storybook/react";

import Code from "../text";

export default {
  title: "Core/Text/Code",
  component: Code,
} as Meta;

const Template = ({ value }) => <Code>{value}</Code>;

export const Primary = Template.bind({});
Primary.args = {
  value: "{key1: [0, 1, 2], key2: 'value', key3: {foo: 'bar'}}",
};
