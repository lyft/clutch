import * as React from "react";
import { styled } from "@clutch-sh/core";
import type { Meta } from "@storybook/react";

import TimeseriesChart, { TimeseriesReferenceLineProps } from "../Charts/timeseries";

export default {
  title: "Core/TimeseriesChart",
  component: TimeseriesChart,
} as Meta;

const ChartContainer = styled("div")({
  width: 900,
  height: 200,
});

export const Primary = () => {
  const mockDataSingleLine = [
    {
      lineName: "line1",
      timestamp: 1546300800000,
      value: 5,
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
    {
      lineName: "line2",
      timestamp: 1546300700000,
      value: 5,
    },
    {
      lineName: "line2",
      timestamp: 1546300600000,
      value: 20,
    },
    {
      lineName: "line2",
      timestamp: 1546301800000,
      value: 30,
    },
  ];

  const mockDataMultiLine = [
    {
      lineName: "line1",
      timestamp: 1546301800,
      value2: 15,
    },
    {
      lineName: "line1",
      timestamp: 1546301900,
      value2: 20,
    },
    {
      lineName: "line1",
      timestamp: 1546302000,
      value2: 80,
    },
    {
      lineName: "line2",
      timestamp: 1546301820,
      value: 10,
    },
    {
      lineName: "line2",
      timestamp: 1546301930,
      value: 20,
    },
    {
      lineName: "line2",
      timestamp: 1546302040,
      value: 40,
    },
  ];

  const mockData3 = [
    {
      lineName: "line1",
      timestamp: 1546301000000,
      value: 10,
    },
    {
      lineName: "line1",
      timestamp: 1546300900000,
      value: 25,
    },
    {
      lineName: "line1",
      timestamp: 1546300800000,
      value: 30,
    },
  ];

  const mockLines = [
    {
      dataKey: "value",
      color: "red",
      animationDuration: 0,
    },
  ];

  const mockLines2 = [
    {
      dataKey: "value",
      color: "purple",
      animationDuration: 2000,
    },
  ];

  const mockMultipleLines = [
    {
      dataKey: "value",
      color: "red",
      animationDuration: 0,
    },
    {
      dataKey: "value2",
      color: "blue",
    },
  ];

  const mockRefLines: TimeseriesReferenceLineProps[] = [
    {
      axis: "x",
      coordinate: 1546300850000,
      label: "ref line 1",
      color: "green",
    },
    {
      axis: "y",
      coordinate: 17,
      label: "ref line 2",
      color: "red",
    },
  ];

  return (
    <ChartContainer>
      <TimeseriesChart
        data={mockDataSingleLine}
        xAxisDataKey="timestamp"
        yAxisDataKey="value"
        lines={mockLines}
      />
      <TimeseriesChart
        data={mockDataMultiLine}
        xAxisDataKey="timestamp"
        yAxisDataKey="value"
        lines={mockMultipleLines}
        singleLineMode={false}
      />
      <TimeseriesChart
        data={mockData3}
        xAxisDataKey="timestamp"
        yAxisDataKey="value"
        lines={mockLines2}
        refLines={mockRefLines}
        drawDots={false}
        enableLegend={false}
      />
    </ChartContainer>
  );
};
