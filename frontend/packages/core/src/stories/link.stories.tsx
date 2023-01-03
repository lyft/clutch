import * as React from "react";
import type { Meta } from "@storybook/react";

import type { LinkProps } from "../link";
import { Link as LinkComponent } from "../link";

export default {
  title: "Core/Link",
  component: LinkComponent,
} as Meta;

const Template = (props: LinkProps) => <LinkComponent {...props}>Clutch Homepage</LinkComponent>;

export const Link = Template.bind({});
Link.args = {
  href: "https://www.clutch.sh",
};
