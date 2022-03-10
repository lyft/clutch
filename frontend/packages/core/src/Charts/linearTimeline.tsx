import React from "react";
import {
  CartesianGrid,
  Legend,
  ScatterChart,
  Scatter,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { calculateTicks, calculateDomainEdges, localTimeFormatter } from "./helpers";

export interface LinearTimelineDataPoint {
    timestamp: number;
    // more to come
  }

  // Note that shape can be a set of shapes denoted by preset strings ("circle", "square", etc.)
  // or a custom SVG element
  export interface LinearTimelineDataPoints {
    points: LinearTimelineDataPoint[];
    shape?: any;
    color?: string;
  }
  
  export interface LinearTimelineData {
    [lane: string]: LinearTimelineDataPoints;
    // ...
  }

  export interface LinearTimelineProps {
    data: LinearTimelineData;
    xAxisDataKey: string;
    regularIntervalTicks?: boolean;
    tickFormatterFunc?: (timestamp: number) => string;
    enableLegend?: boolean;

  }

const LinearTimeline = ({
  data,
  xAxisDataKey = "timestamp",
  regularIntervalTicks = true,
  tickFormatterFunc = localTimeFormatter,
  enableLegend = true,
}: LinearTimelineProps) => {
  const combinedData = Object.keys(data).reduce((acc, lane) => {;
    return [...acc, ...data[lane].points];
  }, []);
  const [xAxisDomainMin, xAxisDomainMax] = calculateDomainEdges(
    combinedData,
    xAxisDataKey,
    .2
  );
  let ticks = [];
  if (regularIntervalTicks) {
    ticks = calculateTicks(combinedData, xAxisDataKey);
  }


  const dataWithIds = Object.keys(data).map((key, index) => {
    const thePoints = data[key].points;
    const pointsWithId = thePoints.map(point => {
      return {
        ...point,
        laneID: index,
      };
    });
    data[key].points = pointsWithId;
    return {
      ...data[key],
      laneID: index,
    };
  })

  return (
    <ResponsiveContainer width="100%" height="100%">
      <ScatterChart >
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis
          dataKey={xAxisDataKey}
          type="number"
          domain={[xAxisDomainMin, xAxisDomainMax]}
          tickFormatter={tickFormatterFunc || null}
          allowDataOverflow
          ticks={regularIntervalTicks ? ticks : null}
        />
        {/* Note due to https://github.com/recharts/recharts/issues/2563 we cannot use a "category" type scatter plot
            To get around this we do a workaround of numbering each lane and hiding the numbers from the user */}
        <YAxis dataKey={"laneID"} type="number" padding={{bottom: 30, top: 30}} hide={true} />
        <Tooltip />
        {enableLegend ? <Legend /> : null}
        {Object.keys(dataWithIds).map((lane) => {
          const points = dataWithIds[lane].points;
          console.log(points);
          return (
            <Scatter
              key={lane}
              data={points}
              name={lane}
              shape={dataWithIds[lane].shape ?? "circle"}
              fill={dataWithIds[lane].color ?? "null"}
            />
          );
        })}
      </ScatterChart>
    </ResponsiveContainer>
  );
};

export default LinearTimeline;
