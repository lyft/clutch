import React from "react";
import type { Meta } from "@storybook/react";

import { TrendingUpIcon } from "../icon";

export default {
  title: "Core/Icon/TrendingUp",
  component: TrendingUpIcon,
} as Meta;

const Template = () => <TrendingUpIcon />;

export const Primary = Template.bind({});
