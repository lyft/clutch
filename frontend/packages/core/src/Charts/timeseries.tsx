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

import { calculateDomainEdges, calculateTicks, localTimeFormatter } from "./helpers";
import type {
  CustomTooltipProps,
  LineProps,
  TimeseriesReferenceLineProps,
  TimeseriesStylingProps,
} from "./types";

/*
  For reference lines (dashed lines), you can set the `axis` property to "x" or "y" to denote which axis 
  the line is on. These can be useful for showing limits or thresholds on graphs.
  There is already a type called `ReferenceLineProps` in Recharts so we avoid that name
*/
export interface TimeseriesChartProps {
  data: any;
  xAxisDataKey?: string;
  yAxisDataKey?: string;
  lines: LineProps[];
  refLines?: TimeseriesReferenceLineProps[];
  singleLineMode?: boolean; // if false, the Y Axis will be based off the max and min of all data combined
  // The assumption is that multiple lines would want to share the same Y Axis.
  drawDots?: boolean;
  legend?: boolean;
  grid?: boolean;
  tickFormatterFunc?: (timeStamp: number) => string;
  xDomainSpread?: number | null;
  yDomainSpread?: number | null;
  connectNulls?: boolean;
  regularIntervalTicks?: boolean;
  tooltipFormatterFunc?: ({ active, payload }: CustomTooltipProps) => JSX.Element;
  stylingProps?: TimeseriesStylingProps;
}

/*
  For the lines, use the `dataKey` property to denote which data points belong to that line. Make sure that
  all the dataKeys match appropriately (the XAxis datakey should correspond to the data that is graphed along
  the XAxis, and same for Y. Note that we currently set the XAxis to be a time scale, hence the name
  Timeseries Chart).
  
  The data will be internally sorted by the XAxis dataKey when using single line mode. When having multiple lines
  (singleLineMode is false) then the user is responsible for sorting and grouping the data appropriately. 
  If they do not sort and group it, there can be discontinuities in the lines. Also, when using multiple lines,
  the user should pass the yAxisDataKey as the biggest ranged y axis datakey in the data, otherwise data will get chopped off.

  The timestamps are interpreted as unix milliseconds.
*/
// TODO(smonero): add tests for this component
const TimeseriesChart = ({
  data,
  xAxisDataKey = "timestamp",
  yAxisDataKey = "value",
  lines,
  refLines,
  singleLineMode = true,
  drawDots = true,
  legend = true,
  grid = true,
  tickFormatterFunc = localTimeFormatter,
  xDomainSpread = 0.2,
  yDomainSpread = 0.2,
  connectNulls = false,
  regularIntervalTicks = false,
  tooltipFormatterFunc = null,
  stylingProps = {},
}: TimeseriesChartProps) => {
  if (singleLineMode) {
    data.sort((a, b) => a[xAxisDataKey] - b[xAxisDataKey]);
  }
  const [yAxisDomainMin, yAxisDomainMax] = calculateDomainEdges(
    data,
    yAxisDataKey,
    yDomainSpread ?? 0
  );
  const [xAxisDomainMin, xAxisDomainMax] = calculateDomainEdges(
    data,
    xAxisDataKey,
    xDomainSpread ?? 0
  );

  // In the spirit of making a friendly UX, there is an option to generate evenly spaced, round-timestamped, ticks.
  // Depending on the distance between the max and min timestamps, we calculate a set of ticks at certain intervals
  // (I.e. if there are several minutes, we might use minute intervals, whereas if there are only a few minutes range,
  // we would use 30 second intervals, and if the range consists of hours, we might use 15 or 30 minute intervals)
  let ticks = [];
  if (regularIntervalTicks) {
    ticks = calculateTicks(data, xAxisDataKey);
  }

  return (
    <ResponsiveContainer width="100%" height="100%">
      <LineChart data={data}>
        {grid ? (
          <CartesianGrid
            stroke={stylingProps?.gridStroke}
            fill={stylingProps?.gridBackgroundColor}
          />
        ) : null}
        <XAxis
          dataKey={xAxisDataKey}
          type="number"
          domain={[xAxisDomainMin, xAxisDomainMax]}
          tickFormatter={tickFormatterFunc}
          allowDataOverflow
          ticks={regularIntervalTicks ? ticks : null}
          stroke={stylingProps?.xAxisStroke}
        />
        {
          // Note that if a number is NaN Recharts will default the domain to `auto`
          singleLineMode ? (
            <YAxis
              dataKey={yAxisDataKey}
              domain={[yAxisDomainMin, yAxisDomainMax]}
              stroke={stylingProps?.yAxisStroke}
              type="number"
            />
          ) : (
            <YAxis
              type="number"
              domain={[yAxisDomainMin, yAxisDomainMax]}
              stroke={stylingProps?.yAxisStroke}
            />
          )
        }
        {/* TODO(smonero): add a default for tooltip formatting */}
        <Tooltip formatter={tooltipFormatterFunc} />
        {legend ? <Legend /> : null}
        {lines
          ? lines.map((line, index) => {
              return (
                <Line
                  key={index.toString() + line.dataKey}
                  type="linear"
                  dataKey={line.dataKey}
                  stroke={line.color}
                  animationDuration={line.animationDuration ?? 0}
                  dot={drawDots}
                  connectNulls={connectNulls}
                />
              );
            })
          : null}
        {refLines &&
          refLines.map((refLine, index) => {
            const props = {};
            props[refLine.axis] = refLine.coordinate;
            return (
              <ReferenceLine
                key={index.toString() + refLine.coordinate.toString()}
                label={refLine.label}
                stroke={refLine.color}
                strokeDasharray={refLine.axis === "x" ? "3 3" : "4 4"}
                {...props}
              />
            );
          })}
      </LineChart>
    </ResponsiveContainer>
  );
};

export default TimeseriesChart;
