import React from "react";
import {
  CartesianGrid,
  Legend,
  Line,
  LineChart,
  ReferenceLine,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

export interface ReferenceLineProps {
  axis: "x" | "y";
  coordinate: number;
  label?: string;
}
export interface LineProps {
  dataKey: string;
  color: string;
  strokeWidth?: number;
  animationDuration?: number;
}
export interface TimeseriesChartProps {
  data: any;
  xAxisDataKey?: string;
  yAxisDataKey?: string;
  lines: LineProps[];
  refLines?: ReferenceLineProps[];
  // TODO: add ref dots, ref areas, zoom enabled, auto colors, legend enabled, cartesian grid options
}

/*
  For the lines, use the `dataKey` property to denote which data points belong to that line. Make sure that
  all the dataKeys match appropriately (the XAxis datakey should correspond to the data that is graphed along
  the XAxis, and same for Y. Note that we currently set the XAxis to be a time scale, hence the name
  Timeseries Chart).
  
  *** VERY IMPORTANT ***
  It is very important that the data is sorted in ascending order by the XAxis. Otherwise the lines will go
  backwards and have problems.
  *** END VERY IMPORTANT ***

  Suggested data format:
  {
    lineName: string
    timestamp: number
    value: number
  }
  For reference lines, you can set the `axis` property to "x" or "y" to denote which axis the line is on.
*/
const TimeseriesChart = ({
  data,
  xAxisDataKey,
  yAxisDataKey,
  lines,
  refLines,
}: TimeseriesChartProps) => {
  return (
    <ResponsiveContainer width="100%" height="100%">
      <LineChart data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis
          dataKey={xAxisDataKey}
          type="number"
          scale="time"
          domain={["dataMin - 1000", "dataMax + 1000"]}
        />
        <YAxis dataKey={yAxisDataKey} domain={["dataMin", "dataMax"]} />
        <Tooltip />
        <Legend />
        {lines
          ? lines.map((line, index) => {
              return (
                <Line
                  key={index.toString() + line.dataKey + line.color}
                  type="linear"
                  dataKey={line.dataKey}
                  stroke={line.color}
                  animationDuration={line.animationDuration !== null ? line.animationDuration : 500}
                />
              );
            })
          : null}
        {refLines
          ? refLines.map((refLine, index) => {
              return refLine.axis === "x" ? (
                <ReferenceLine
                  key={index.toString() + refLine.axis + refLine.coordinate.toString()}
                  x={refLine.coordinate}
                  label={refLine.label}
                />
              ) : (
                <ReferenceLine
                  key={index.toString() + refLine.axis + refLine.coordinate.toString()}
                  y={refLine.coordinate}
                  label={refLine.label}
                />
              );
            })
          : null}
      </LineChart>
    </ResponsiveContainer>
  );
};

export default TimeseriesChart;
