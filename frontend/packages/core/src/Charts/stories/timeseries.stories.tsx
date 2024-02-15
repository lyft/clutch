import React from "react";
import { alpha } from "@mui/system";
import type { Meta } from "@storybook/react";

import { useTheme } from "../../AppProvider/themes";
import styled from "../../styled";
import { dateTimeFormatter, isoTimeFormatter } from "../helpers";
import TimeseriesChart from "../timeseries";
import type { TimeseriesReferenceLineProps } from "../types";

export default {
  title: "Core/Charts/TimeseriesChart",
  component: TimeseriesChart,
} as Meta;

const ChartContainer = styled("div")({
  width: 900,
  height: 400,
});

export const SingleDataLine = () => {
  const theme = useTheme();
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
      lineName: "line1",
      timestamp: 1546300700000,
      value: 5,
    },
    {
      lineName: "line1",
      timestamp: 1546300600000,
      value: 20,
    },
    {
      lineName: "line1",
      timestamp: 1546301800000,
      value: 30,
    },
  ];

  const mockLines = [
    {
      dataKey: "value",
      color: theme.colors.red[500],
      animationDuration: 0,
    },
  ];

  return (
    <ChartContainer>
      <TimeseriesChart
        data={mockDataSingleLine}
        xAxisDataKey="timestamp"
        yAxisDataKey="value"
        lines={mockLines}
        xDomainSpread={0.3}
        yDomainSpread={0.3}
        regularIntervalTicks
      />
      <TimeseriesChart
        data={mockDataSingleLine}
        xAxisDataKey="timestamp"
        yAxisDataKey="value"
        lines={mockLines}
        xDomainSpread={0.3}
        yDomainSpread={0.3}
        regularIntervalTicks
        tickFormatterFunc={isoTimeFormatter}
      />
      <TimeseriesChart
        data={mockDataSingleLine}
        xAxisDataKey="timestamp"
        yAxisDataKey="value"
        lines={mockLines}
        xDomainSpread={0.3}
        yDomainSpread={0.3}
        regularIntervalTicks
        tickFormatterFunc={dateTimeFormatter}
      />
    </ChartContainer>
  );
};

export const MultiLine = () => {
  const theme = useTheme();
  const mockDataMultiLine = [
    {
      lineName: "line1",
      timestamp: 1546201800,
      value2: 15,
    },
    {
      lineName: "line1",
      timestamp: 1546291900,
      value2: 20,
    },
    {
      lineName: "line1",
      timestamp: 1546302000,
      value2: 80,
    },
    {
      lineName: "line2",
      timestamp: 1546201820,
      value: 10,
    },
    {
      lineName: "line2",
      timestamp: 1546291930,
      value: 20,
    },
    {
      lineName: "line2",
      timestamp: 1546302040,
      value: 40,
    },
  ];

  const mockMultipleLines = [
    {
      dataKey: "value",
      color: theme.colors.red[500],
      animationDuration: 0,
    },
    {
      dataKey: "value2",
      color: theme.colors.blue[500],
    },
  ];

  return (
    <ChartContainer>
      <TimeseriesChart
        data={mockDataMultiLine}
        xAxisDataKey="timestamp"
        yAxisDataKey="value2"
        lines={mockMultipleLines}
        singleLineMode={false}
        xDomainSpread={0.3}
        yDomainSpread={0.3}
      />
    </ChartContainer>
  );
};

/** *
 * This example shows that users can have raw values rather than using a formatter func
 * for the ticks along the X-Axis. It also shows reference lines and the disabling of
 * other options.
 */
export const ReferenceLinesAndThingsDisabled = () => {
  const theme = useTheme();
  const mockData = [
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
      color: theme.colors.charts.common.data[0],
      animationDuration: 2000,
    },
  ];

  const mockRefLines: TimeseriesReferenceLineProps[] = [
    {
      axis: "x",
      coordinate: 1546300850000,
      label: "ref line 1",
      color: theme.colors.green[500],
    },
    {
      axis: "y",
      coordinate: 17,
      label: "ref line 2",
      color: theme.colors.red[500],
    },
  ];

  return (
    <ChartContainer>
      <TimeseriesChart
        data={mockData}
        xAxisDataKey="timestamp"
        yAxisDataKey="value"
        lines={mockLines}
        refLines={mockRefLines}
        drawDots={false}
        legend={false}
        grid={false}
        tickFormatterFunc={v => `${v}`}
        xDomainSpread={null}
      />
    </ChartContainer>
  );
};

export const StyledChart = () => {
  const theme = useTheme();
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
      lineName: "line1",
      timestamp: 1546300700000,
      value: 5,
    },
    {
      lineName: "line1",
      timestamp: 1546300600000,
      value: 20,
    },
    {
      lineName: "line1",
      timestamp: 1546301800000,
      value: 30,
    },
  ];

  const mockLines = [
    {
      dataKey: "value",
      color: theme.colors.red[500],
      animationDuration: 0,
    },
  ];

  return (
    <ChartContainer>
      <TimeseriesChart
        data={mockDataSingleLine}
        xAxisDataKey="timestamp"
        yAxisDataKey="value"
        lines={mockLines}
        xDomainSpread={0.3}
        yDomainSpread={0.3}
        regularIntervalTicks
        stylingProps={{
          gridBackgroundColor: alpha(theme.colors.red[600], 0.4),
          gridStroke: theme.colors.blue[700],
          xAxisStroke: theme.colors.green[500],
          yAxisStroke: theme.colors.amber[500],
        }}
      />
    </ChartContainer>
  );
};
