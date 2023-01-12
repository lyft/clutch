import * as React from "react";
import type { Meta } from "@storybook/react";

import CodeComponent from "../text";

export default {
  title: "Core/Text/Code",
  component: CodeComponent,
} as Meta;

const Template = ({ value }) => <CodeComponent>{value}</CodeComponent>;

export const Code = Template.bind({});
Code.args = {
  value: "{key1: [0, 1, 2], key2: 'value', key3: {foo: 'bar'}}",
};
