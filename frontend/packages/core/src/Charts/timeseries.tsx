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

import { calculateDomainEdges, localTimeFormatter } from "./helpers";

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
  drawDots?: boolean;
  enableLegend?: boolean;
  enableGrid?: boolean;
  tickFormatterFunc?: (timeStamp: number) => string;
  xDomainSpread?: number | null;
  yDomainSpread?: number | null;
  connectNulls?: boolean;
  friendlyTicks?: boolean;
  // TODO: add ref dots, ref areas, zoom enabled, auto colors,
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
  If they do not sort and group it, there can be discontinuities in the lines. Also, when using multiple lines,
  the user should pass the yAxisDataKey as the biggest ranged y axis datakey in the data, otherwise data will get chopped off.
  *** END VERY IMPORTANT ***
*/
const TimeseriesChart = ({
  data,
  xAxisDataKey,
  yAxisDataKey,
  lines,
  refLines,
  singleLineMode = true,
  drawDots = true,
  enableLegend = true,
  enableGrid = true,
  tickFormatterFunc = localTimeFormatter,
  xDomainSpread,
  yDomainSpread,
  connectNulls = false,
  friendlyTicks = false,
}: TimeseriesChartProps) => {
  if (singleLineMode) {
    data.sort((a, b) => a[xAxisDataKey] - b[xAxisDataKey]);
  }
  const [yAxisDomainMin, yAxisDomainMax] = calculateDomainEdges(
    data,
    yAxisDataKey,
    yDomainSpread === null ? 0 : yDomainSpread
  );
  const [xAxisDomainMin, xAxisDomainMax] = calculateDomainEdges(
    data,
    xAxisDataKey,
    xDomainSpread === null ? 0 : xDomainSpread
  );

  // In the spirit of making a friendly UX, there is an option to generate evenly spaced, round-timestamped, ticks.
  // Depending on the distance between the max and min timestamps, we calculate a set of ticks at certain intervals
  // (I.e. if there are several minutes, we might use minute intervals, whereas if there are only a few minutes range,
  // we would use 30 second intervals, and if the range consists of hours, we might use 15 or 30 minute intervals)
  if (friendlyTicks) {
  }

  return (
    <ResponsiveContainer width="100%" height="100%">
      <LineChart data={data}>
        {enableGrid ? <CartesianGrid /> : null}
        <XAxis
          dataKey={xAxisDataKey}
          type="number"
          scale="linear"
          domain={[xAxisDomainMin, xAxisDomainMax]}
          tickFormatter={tickFormatterFunc}
          allowDataOverflow={true}
        />
        {
          // Note that if a number is NaN Recharts will default the domain to `auto`
          singleLineMode ? (
            <YAxis dataKey={yAxisDataKey} domain={[yAxisDomainMin, yAxisDomainMax]} type="number" />
          ) : (
            <YAxis type="number" domain={[yAxisDomainMin, yAxisDomainMax]} />
          )
        }
        <Tooltip labelFormatter={tickFormatterFunc} />
        {enableLegend ? <Legend /> : null}
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
                  dot={drawDots}
                  connectNulls={connectNulls}
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
                  strokeDasharray="4 4"
                />
              );
            })
          : null}
      </LineChart>
    </ResponsiveContainer>
  );
};

export default TimeseriesChart;
