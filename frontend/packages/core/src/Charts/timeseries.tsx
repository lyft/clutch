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

const tickFormatterFunc = timeStamp => {
  const date = new Date(timeStamp);
  return date.toLocaleTimeString();
};

export type ReferenceLineAxis = "x" | "y";
/*
  For reference lines (dashed lines), you can set the `axis` property to "x" or "y" to denote which axis 
  the line is on. These can be useful for showing limits or thresholds on graphs.
  There is already a type called `ReferenceLineProps` in Recharts so we avoid that name
*/
export interface TimeseriesReferenceLineProps {
  axis: ReferenceLineAxis;
  coordinate: number;
  color: string;
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
  refLines?: TimeseriesReferenceLineProps[];
  singleLineMode?: boolean; // if false, the Y Axis will be based off the max and min of all data combined
  // The assumption is that multiple lines would want to share the same Y Axis.
  // TODO: add ref dots, ref areas, zoom enabled, auto colors, legend enabled,
  // tooltip options, activeDot options, dark mode / styling,
}

/*
  For the lines, use the `dataKey` property to denote which data points belong to that line. Make sure that
  all the dataKeys match appropriately (the XAxis datakey should correspond to the data that is graphed along
  the XAxis, and same for Y. Note that we currently set the XAxis to be a time scale, hence the name
  Timeseries Chart).
  
  *** VERY IMPORTANT ***
  The data will be internally sorted by the XAxis dataKey when using single line mode. When having multiple lines
  (singleLineMode is false) then the user is responsible for sorting and grouping the data appropriately. 
  If they do not sort and group it, there can be discontinuities in the lines.
  *** END VERY IMPORTANT ***
*/
const TimeseriesChart = ({
  data,
  xAxisDataKey,
  yAxisDataKey,
  lines,
  refLines,
  singleLineMode = true,
}: TimeseriesChartProps) => {
  if (singleLineMode) {
    data.sort((a, b) => a[xAxisDataKey] - b[xAxisDataKey]);
  }
  return (
    <ResponsiveContainer width="100%" height="100%">
      <LineChart data={data}>
        <CartesianGrid />
        <XAxis
          dataKey={xAxisDataKey}
          type="number"
          scale="time"
          domain={["dataMin - 1000", "dataMax + 1000"]}
          tickFormatter={tickFormatterFunc}
        />
        {singleLineMode ? (
          <YAxis dataKey={yAxisDataKey} domain={["dataMin", "dataMax"]} type="number" />
        ) : (
          <YAxis type="number" />
        )}
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
                  animationDuration={
                    line.animationDuration !== null ? line.animationDuration : null
                  }
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
                  stroke={refLine.color}
                  strokeDasharray="3 3"
                />
              ) : (
                <ReferenceLine
                  key={index.toString() + refLine.axis + refLine.coordinate.toString()}
                  y={refLine.coordinate}
                  label={refLine.label}
                  stroke={refLine.color}
                  strokeDasharray="3 4"
                />
              );
            })
          : null}
      </LineChart>
    </ResponsiveContainer>
  );
};

export default TimeseriesChart;
