import type { ReactElement } from "react";

export type ReferenceLineAxis = "x" | "y";
export interface LinearTimelineDataPoint {
  timestamp: number;
}

export type PresetShape = "circle" | "square" | "cross" | "diamond" | "star" | "triangle" | "wye";
// Note that shape can be a preset of one of the shapes or a custom SVG element
// See https://recharts.org/en-US/api/Scatter#shape for more details
export interface LinearTimelineDataPoints {
  points: LinearTimelineDataPoint[];
  shape?: ReactElement<SVGElement> | ((props: any) => ReactElement<SVGElement>) | PresetShape;
  color?: string;
}

export interface LinearTimelineData {
  [lane: string]: LinearTimelineDataPoints;
}

export interface LinearTimelineStylingProps {
  xAxisStroke?: string;
  // note no yAxis stroke because y Axis is hidden
  tooltipBackgroundColor?: string;
  tooltipTextColor?: string;
  gridBackgroundColor?: string;
  gridStroke?: string;
  // TODO(smonero): add size control of icons via z-axis
}

export interface CustomTooltipProps {
  active: boolean;
  payload: any; // A huge object that contains all the info for the data point and more
}

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

export interface TimeseriesStylingProps {
  xAxisStroke?: string;
  yAxisStroke?: string;
  gridBackgroundColor?: string;
  gridStroke?: string;
}

export interface PieChartData {
  name: string;
  value: string | number;
  color?: string;
}
