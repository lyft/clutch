import React from "react";
import {
  CartesianGrid,
  Legend,
  ResponsiveContainer,
  Scatter,
  ScatterChart,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

import { calculateDomainEdges, calculateTicks, localTimeFormatter } from "./helpers";

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
  legend?: boolean;
}

const LinearTimeline = ({
  data,
  xAxisDataKey = "timestamp",
  regularIntervalTicks = true,
  tickFormatterFunc = localTimeFormatter,
  legend = true,
}: LinearTimelineProps) => {
  const combinedData = Object.keys(data).reduce((acc, lane) => {
    return [...acc, ...data[lane].points];
  }, []);
  const [xAxisDomainMin, xAxisDomainMax] = calculateDomainEdges(combinedData, xAxisDataKey, 0.2);
  let ticks = [];
  if (regularIntervalTicks) {
    ticks = calculateTicks(combinedData, xAxisDataKey);
  }

  // Because we can't rely on using "category" for the Y-Axis, we need to assign Ids (based off the index)
  const mappingOfLaneIdsToNames = {};
  const dataWithIds = Object.keys(data).map((key, index) => {
    mappingOfLaneIdsToNames[index] = key;
    const thePoints = data[key].points;
    const pointsWithId = thePoints.map(point => {
      return {
        ...point,
        laneID: index,
      };
    });
    return {
      points: pointsWithId,
      shape: data[key].shape,
      color: data[key].color,
      laneID: index,
    };
  });
  const formatIdsToNames = (value: string) => {
    return <span>{mappingOfLaneIdsToNames[value]}</span>;
  };

  // TODO: Allow for proper styling and make things less hacky than "payload[0]"
  const formatTooltip = ({ active, payload, label }) => {
    if (active) {
      return <div>{localTimeFormatter(payload[0].value)}</div>;
    }

    return null;
  };

  return (
    <ResponsiveContainer width="100%" height="100%">
      <ScatterChart>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis
          dataKey={xAxisDataKey}
          type="number"
          domain={[xAxisDomainMin, xAxisDomainMax]}
          tickFormatter={tickFormatterFunc}
          allowDataOverflow
          ticks={regularIntervalTicks ? ticks : null}
        />
        {/* Note due to https://github.com/recharts/recharts/issues/2563 we cannot use a "category" type scatter plot
            To get around this we do a workaround of numbering each lane and hiding the numbers from the user */}
        <YAxis dataKey="laneID" type="number" padding={{ bottom: 30, top: 30 }} hide />
        <Tooltip content={formatTooltip} />
        {/* TODO: Use the Z-Axis for a "zoom" amount to enlarge or shrink icon size */}
        {legend ? <Legend iconSize={18} formatter={formatIdsToNames} /> : null}
        {Object.keys(dataWithIds).map(lane => {
          const { points } = dataWithIds[lane];
          return (
            <Scatter
              key={lane}
              data={points}
              name={lane}
              shape={dataWithIds[lane].shape ?? "circle"}
              fill={dataWithIds[lane].color ?? "null"}
              legendType={dataWithIds[lane].shape ?? "circle"}
            />
          );
        })}
      </ScatterChart>
    </ResponsiveContainer>
  );
};

export default LinearTimeline;
