import * as React from "react";
import type { Meta } from "@storybook/react";

import type { CardProps } from "../card";
import { Card as CardComponent, CardContent } from "../card";

export default {
  title: "Core/Card/Basic",
  component: CardComponent,
} as Meta;

const Template = (props: CardProps) => <CardComponent {...props} />;

export const Basic = Template.bind({});
Basic.args = {
  children: <CardContent>Hello world!</CardContent>,
};
