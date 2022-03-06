import React from "react";
import {
  CartesianGrid,
  Legend,
  Line,
  ScatterChart,
  Scatter,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

export interface LineProps {
  dataKey: string;
  color: string;
  strokeWidth?: number;
  animationDuration?: number;
}
export interface LinearTimelineProps {
  data: any;
  xAxisDataKey?: string;
  yAxisDataKey?: string;
  lines: LineProps[];
  // TODO: add ref dots, ref areas, zoom enabled, auto colors, legend enabled, cartesian grid options
}

/*

  
*/
interface LinearTimelineDataPoint {
    timestamp: number | string | Date | Long;
    lane: string;
    metadata: any;
    // more to come
  }
  
  interface LInearTimelineData {
    [lane: string]: LinearTimelineDataPoint[];
    // ...
  }
  

const LinearTimeline = ({
  data,
  xAxisDataKey,
  yAxisDataKey,
  yAxisLaneIds,
}: LinearTimelineProps) => {


  return (
    <ResponsiveContainer width="100%" height="100%">
      <ScatterChart data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis
          dataKey={xAxisDataKey}
          type="number"
          scale="time"
          domain={["dataMin - 1000", "dataMax + 1000"]}
        />
        <YAxis dataKey={"lane"} type="category" />
        <Tooltip />
        <Legend />
        {data
          ? data.map((d, index) => {
              return (
                <Scatter
                  key={index.toString() + d[yAxisDataKey}
                  data={data[d]}
                  shape={d.shape ? d.shape : "circle"}
                />
              )
            })
          : null}
      </ScatterChart>
    </ResponsiveContainer>
  );
};

export default LinearTimeline;
