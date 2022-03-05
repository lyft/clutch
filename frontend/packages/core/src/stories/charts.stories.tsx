import * as React from "react";
import { styled } from "@clutch-sh/core";
import type { Meta } from "@storybook/react";

import TimeseriesChart from "../Charts/timeseries-chart";

export default {
  title: "Core/TimeseriesChart",
  component: TimeseriesChart,
} as Meta;

const ChartContainer = styled("div")({
  width: 600,
  height: 400,
});

export const Primary = () => {
  const mockData = [
    {
      lineName: "line1",
      timestamp: 1546300800000,
      value: 10,
    },
    {
      lineName: "line1",
      timestamp: 1546300900000,
      value: 20,
    },
    {
      lineName: "line1",
      timestamp: 1546301000000,
      value: 30,
    },
  ];

  const mockLines = [
    {
      dataKey: "value",
      color: "red",
    },
  ];
  // data
  // reference lines
  // TimeseriesChart
  // const TimeseriesChart = ({data, xAxisDataKey, yAxisDataKey, lines, refLines }: TimeseriesChartProps) => {
  return (
    <ChartContainer>
      <TimeseriesChart
        data={mockData}
        xAxisDataKey="timestamp"
        yAxisDataKey="value"
        lines={mockLines}
      />
    </ChartContainer>
  );
};
