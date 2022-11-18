import * as React from "react";
import type { Meta } from "@storybook/react";

import TimeAgo from "../timeago";

export default {
  title: "Core/TimeAgo",
  component: TimeAgo,
} as Meta;

const Template = ({ date, short }) => <TimeAgo date={date} short={short} />;

export const Default = Template.bind({});

Default.args = {
  date: 1668788531 * 1000,
  short: false,
};
