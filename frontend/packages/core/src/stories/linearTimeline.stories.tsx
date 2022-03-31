import * as React from "react";
import styled from "@emotion/styled";
import type { Meta } from "@storybook/react";

import LinearTimeline from "../Charts/linearTimeline";
import type { LinearTimelineData } from "../Charts/types";

export default {
  title: "Core/Charts/LinearTimeline",
  component: LinearTimeline,
} as Meta;

const ChartContainer = styled("div")({
  width: 400,
  height: 200,
});

export const Primary = () => {
  const mockData: LinearTimelineData = {
    deploys: { points: [{ timestamp: 1588884888 }], shape: "cross", color: "purple" },
    "k8s events": { points: [{ timestamp: 1588985888 }], shape: "triangle", color: "blue" },
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
      color: "green",
    },
  };
  return (
    <ChartContainer>
      <LinearTimeline data={mockData} xAxisDataKey="timestamp" />
    </ChartContainer>
  );
};

export const ColoredChart = () => {
  const mockData: LinearTimelineData = {
    deploys: { points: [{ timestamp: 1588884888 }], shape: "cross", color: "purple" },
    "k8s events": { points: [{ timestamp: 1588985888 }], shape: "triangle", color: "blue" },
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
      color: "black",
    },
  };
  return (
    <ChartContainer>
      <LinearTimeline
        data={mockData}
        xAxisDataKey="timestamp"
        stylingProps={{
          xAxisStroke: "red",
          tooltipBackgroundColor: "blue",
          tooltipTextColor: "white",
          gridBackgroundColor: "green",
          gridStroke: "yellow",
        }}
      />
    </ChartContainer>
  );
};
