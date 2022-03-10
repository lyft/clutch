import * as React from "react";
import { styled } from "@clutch-sh/core";
import type { Meta } from "@storybook/react";

import LinearTimeline from "../Charts/linearTimeline"
import type { LinearTimelineData } from "../Charts/linearTimeline";

export default {
  title: "Core/LinearTimeline",
  component: LinearTimeline,
} as Meta;

const ChartContainer = styled("div")({
  width: 400,
  height: 200,
});

export const Primary = () => {
  const mockData: LinearTimelineData = {
    "deploys": { points: [ { timestamp: 1588884888 } ], shape: "cross", color: "purple" },
    "k8s events": { points: [ { timestamp: 1588985888 } ], shape: "wye", color: "pink" },
    "explosions": { points: [ { timestamp: 1588788888 }, { timestamp: 1589708888 }, 
      { timestamp: 1589608088 }, 
      { timestamp: 1589618088 },
      { timestamp: 1589628088 },
      { timestamp: 1589638088 },
      { timestamp: 1589648088 },
      { timestamp: 1589658088 },
      { timestamp: 1589668088 },
      { timestamp: 1589678088 },
      { timestamp: 1589688088 },
      { timestamp: 1589698088 },
      { timestamp: 1589607088 },
    ], shape: "star", color: "green" },
  };
  return (
    <ChartContainer>
      <LinearTimeline
        data={mockData}
        xAxisDataKey="timestamp"
      />
    </ChartContainer>
  );
};
