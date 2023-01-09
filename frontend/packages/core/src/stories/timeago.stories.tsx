import * as React from "react";
import type { Meta } from "@storybook/react";

import TimeAgoComponent from "../timeago";

export default {
  title: "Core/TimeAgo",
  component: TimeAgoComponent,
} as Meta;

const Template = ({ date, live, short }) => (
  <TimeAgoComponent date={date} live={live} short={short} />
);

export const TimeAgo = Template.bind({});

TimeAgo.args = {
  date: 1668788531 * 1000,
  short: false,
  live: false,
};
