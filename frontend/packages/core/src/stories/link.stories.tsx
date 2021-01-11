import * as React from "react";
import type { Meta } from "@storybook/react";

import type { LinkProps } from "../link";
import { Link } from "../link";

export default {
  title: "Core/Link",
  component: Link,
} as Meta;

const Template = (props: LinkProps) => <Link {...props}>Clutch Homepage</Link>;

export const Default = Template.bind({});
Default.args = {
  href: "https://www.clutch.sh",
};
