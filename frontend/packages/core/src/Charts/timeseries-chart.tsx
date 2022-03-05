import React, { useEffect } from "react";
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, ReferenceLine } from 'recharts';

  
export interface ReferenceLineProps {
  axis: "x" | "y";
  coordinate: number;
  label?: string;
}
export interface LineProps {
    dataKey: string;
    color: string;
}
export interface TimeseriesChartProps {
  data: any;
  xAxisDataKey?: string;
  yAxisDataKey?: string;
  lines: LineProps[];
  refLines?: ReferenceLineProps[];
  // To add: ref dots, ref areas, zoom enabled, auto colors, legend enabled, cartesian grid options
};

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
const TimeseriesChart = ({data, xAxisDataKey, yAxisDataKey, lines, refLines }: TimeseriesChartProps) => {
    return (
        <ResponsiveContainer width="100%" height="100%">
          <LineChart
            data={data}
          >
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey={xAxisDataKey} type="number" scale="time" />
            <YAxis dataKey={yAxisDataKey} />
            <Tooltip />
            <Legend />
            {
              lines.map((line, index) => {
                return (
                  <Line key={index} type="linear" dataKey={line.dataKey} stroke={line.color} />
                )
              })
            }
            {
              refLines.map(refLine => {
                return (refLine.axis === "x" ? 
                <ReferenceLine x={refLine.coordinate} label={refLine.label}  /> : <ReferenceLine y={refLine.coordinate} label={refLine.label} />)
              })
            }
          </LineChart>
        </ResponsiveContainer>
      );
  };
  
  export default TimeseriesChart;
  
