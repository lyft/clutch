import * as React from "react";
import type { Meta } from "@storybook/react";

import TimeseriesChart from "../timeseries-chart";

export default {
  title: "Core/TimeseriesChart",
  component: TimeseriesChart,
} as Meta;

export const Primary = () => {
  //data
  // reference lines
  // TimeseriesChart
  return <TimeseriesChart />;
};
