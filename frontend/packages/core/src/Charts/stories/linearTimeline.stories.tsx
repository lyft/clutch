import React from "react";
import type { Meta } from "@storybook/react";

import { useTheme } from "../../AppProvider/themes";
import styled from "../../styled";
import LinearTimeline from "../linearTimeline";
import type { LinearTimelineData } from "../types";

export default {
  title: "Core/Charts/LinearTimeline",
  component: LinearTimeline,
} as Meta;

const ChartContainer = styled("div")({
  width: 400,
  height: 200,
});

export const Primary = () => {
  const theme = useTheme();
  const mockData: LinearTimelineData = {
    deploys: {
      points: [{ timestamp: 1588884888 }],
      shape: "cross",
      color: theme.colors.charts.common.data[1],
    },
    "k8s events": {
      points: [{ timestamp: 1588985888 }],
      shape: "triangle",
      color: theme.colors.charts.common.data[6],
    },
    explosions: {
      points: [
        { timestamp: 1588788888 },
        { timestamp: 1589708888 },
        { timestamp: 1589608088 },
        { timestamp: 1589618088 },
        { timestamp: 1589828088 },
        { timestamp: 1589138088 },
        { timestamp: 1589248088 },
        { timestamp: 1589358088 },
        { timestamp: 1589468088 },
        { timestamp: 1589508088 },
        { timestamp: 1589688088 },
        { timestamp: 1589798088 },
        { timestamp: 1589807088 },
      ],
      shape: "star",
      color: theme.colors.charts.common.data[2],
    },
  };
  return (
    <ChartContainer>
      <LinearTimeline data={mockData} xAxisDataKey="timestamp" />
    </ChartContainer>
  );
};

export const ColoredChart = () => {
  const theme = useTheme();
  const mockData: LinearTimelineData = {
    deploys: {
      points: [{ timestamp: 1588884888 }],
      shape: "cross",
      color: theme.colors.red[500],
    },
    "k8s events": {
      points: [{ timestamp: 1588985888 }],
      shape: "triangle",
      color: theme.colors.amber[500],
    },
    explosions: {
      points: [
        { timestamp: 1588788888 },
        { timestamp: 1589708888 },
        { timestamp: 1589608088 },
        { timestamp: 1589618088 },
        { timestamp: 1589828088 },
        { timestamp: 1589138088 },
        { timestamp: 1589248088 },
        { timestamp: 1589358088 },
        { timestamp: 1589468088 },
        { timestamp: 1589508088 },
        { timestamp: 1589688088 },
        { timestamp: 1589798088 },
        { timestamp: 1589807088 },
      ],
      shape: "star",
      color: theme.colors.green[500],
    },
  };
  return (
    <ChartContainer>
      <LinearTimeline
        data={mockData}
        xAxisDataKey="timestamp"
        stylingProps={{
          xAxisStroke: theme.colors.blue[500],
          tooltipBackgroundColor: theme.colors.blue[200],
          tooltipTextColor: theme.colors.blue[900],
          gridBackgroundColor: theme.colors.blue[100],
          gridStroke: theme.colors.blue[300],
        }}
      />
    </ChartContainer>
  );
};
