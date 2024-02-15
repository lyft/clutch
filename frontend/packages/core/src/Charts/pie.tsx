import React, { PureComponent } from "react";
import { ThemeContext } from "@emotion/react";
import type { Theme } from "@mui/material";
import {
  Cell,
  Label,
  Legend,
  Pie,
  PieChart as RechartsPieChart,
  ResponsiveContainer,
  Sector,
  Tooltip,
} from "recharts";

import styled from "../styled";

import type { PieChartData } from "./types";

export interface PieChartProps {
  /**
   * The data to display in the chart
   */
  data: PieChartData[];
  /**
   * Optional dimensions for the chart
   */
  dimensions?: {
    height?: number;
    width?: number;
    innerRadius?: number;
    outerRadius?: number;
    paddingAngle?: number;
    cx?: string;
    cy?: string;
  };
  /**
   * If `true` will display a label for the pie pieces
   * @default false
   */
  label?: boolean | React.ReactElement;
  /**
   * If `true` will display a line to the label, only active with `label`
   * @default false
   */
  labelLine?: boolean;
  /**
   * If set, will display a label in the center of the pie chart (title) with an optional second label (subtitle)
   */
  centerLabel?: {
    title?: string;
    subtitle?: string;
  };
  /**
   * If `true` will display a legend for the chart
   * @default false
   */
  legend?: boolean;
  /**
   * Settings for responsive container
   * @default true
   */
  responsive?:
    | boolean
    | {
        width?: number | string;
        height?: number | string;
        aspect?: number;
      };
  /**
   * It will display an active tooltip with changing text
   * @default on
   */
  activeTooltip?:
    | boolean
    | {
        staticLabel?: string;
        payloadLabel?: string;
        formatter?: (payload: any) => string;
      };
  /**
   * If `true` will display a tooltip on hover over the chart slice
   * @default false
   */
  tooltip?: boolean;
  /**
   * (Optional) children to render inside of the <PieChart />, can reference API from
   * Recharts (https://recharts.org/en-US/api/PieChart)
   * These can override the Legend and Tooltip provided via the component
   */
  children?: React.ReactChild | React.ReactChild[];
}

interface PieChartState {
  activeIndex?: number;
}

const ChartLabelPrimary = styled("text")(({ theme }: { theme: Theme }) => ({
  fill: theme.colors.charts.pie.labelPrimary,
}));

const ChartLabelSecondary = styled("text")(({ theme }: { theme: Theme }) => ({
  fill: theme.colors.charts.pie.labelSecondary,
}));

const renderActiveShape = (props, options) => {
  const RADIAN = Math.PI / 180;
  const {
    cx,
    cy,
    midAngle,
    innerRadius,
    outerRadius,
    startAngle,
    endAngle,
    fill,
    payload,
    percent,
    value,
  } = props;
  const sin = Math.sin(-RADIAN * midAngle);
  const cos = Math.cos(-RADIAN * midAngle);
  const sx = cx + (outerRadius + 10) * cos;
  const sy = cy + (outerRadius + 10) * sin;
  const mx = cx + (outerRadius + 30) * cos;
  const my = cy + (outerRadius + 30) * sin;
  const ex = mx + (cos >= 0 ? 1 : -1) * 22;
  const ey = my;
  const textAnchor = cos >= 0 ? "start" : "end";

  return (
    <g>
      {(options.formatter || options.staticLabel || options.payloadLabel) && (
        <text x={cx} y={cy} dy={8} textAnchor="middle">
          {options.formatter && options.formatter(payload)}
          {options.staticLabel && options.staticLabel}
          {options.payloadLabel && `${options.payloadLabel} ${payload.name}`}
        </text>
      )}
      <Sector
        cx={cx}
        cy={cy}
        innerRadius={innerRadius}
        outerRadius={outerRadius}
        startAngle={startAngle}
        endAngle={endAngle}
        fill={fill}
      />
      <Sector
        cx={cx}
        cy={cy}
        startAngle={startAngle}
        endAngle={endAngle}
        innerRadius={outerRadius + 6}
        outerRadius={outerRadius + 10}
        fill={fill}
      />
      <path d={`M${sx},${sy}L${mx},${my}L${ex},${ey}`} stroke={fill} fill="none" />
      <circle cx={ex} cy={ey} r={2} fill={fill} stroke="none" />
      <ChartLabelPrimary x={ex + (cos >= 0 ? 1 : -1) * 12} y={ey} textAnchor={textAnchor}>
        {payload.name}
      </ChartLabelPrimary>
      <ChartLabelSecondary x={ex + (cos >= 0 ? 1 : -1) * 12} y={ey} dy={18} textAnchor={textAnchor}>
        {`${value} (${(percent * 100).toFixed(2)}%)`}
      </ChartLabelSecondary>
    </g>
  );
};

const CenterLabel = props => {
  const { options, viewBox } = props;
  const { cx, cy, fill } = viewBox;

  if (!options) {
    return null;
  }

  return (
    <g>
      {options.title && (
        <text x={cx} y={cy} textAnchor="middle" fill={fill} style={{ fontSize: "36px" }}>
          {options.title}
        </text>
      )}
      {options.subtitle && (
        <text x={cx} y={cy} dy={28} textAnchor="middle" fill={fill} style={{ fontSize: "14px" }}>
          {options.subtitle}
        </text>
      )}
    </g>
  );
};

class PieChart extends PureComponent<PieChartProps, PieChartState> {
  constructor(props) {
    super(props);
    this.state = { activeIndex: 0 };
  }

  onPieEnter = (_, activeIndex) => {
    this.setState({ activeIndex });
  };

  render() {
    const {
      children,
      centerLabel,
      data,
      dimensions,
      activeTooltip,
      label,
      labelLine,
      legend = false,
      responsive = true,
      tooltip,
    } = this.props;

    const { colors } = this.context;

    const chartOptions = {
      activeTooltip: typeof activeTooltip === "boolean" ? activeTooltip : true,
      activeTooltipOptions: typeof activeTooltip !== "boolean" ? { ...activeTooltip } : {},
      responsive: typeof responsive === "boolean" ? responsive : true,
      responsiveDimensions: {
        width: "99%",
        height: "99%",
        aspect: 2,
        ...(typeof responsive !== "boolean" ? responsive : {}),
      },
      centerLabel,
      label,
      labelLine,
      legend,
      tooltip,
      dimensions: {
        height: 275,
        width: 275,
        innerRadius: 60,
        outerRadius: 80,
        paddingAngle: 2,
        cx: "50%",
        cy: "50%",
        ...(dimensions || {}),
      },
    };

    const additionalProps = {
      ...(chartOptions.activeTooltip
        ? {
            // eslint-disable-next-line react/destructuring-assignment
            activeIndex: this.state?.activeIndex,
            activeShape: props =>
              renderActiveShape(props, { ...chartOptions.activeTooltipOptions }),
          }
        : {}),
      ...(chartOptions.label
        ? {
            label: chartOptions.label,
            labelLine: chartOptions.labelLine,
          }
        : {}),
    };

    const chart = (
      <RechartsPieChart
        height={chartOptions.dimensions.height}
        width={chartOptions.dimensions.width}
      >
        <Pie
          data={data}
          fill={colors.charts.common.data[0]}
          dataKey="value"
          onMouseEnter={this.onPieEnter}
          {...chartOptions.dimensions}
          {...additionalProps}
        >
          {data.map((entry, index) => (
            <Cell
              // eslint-disable-next-line react/no-array-index-key
              key={`cell-${index}`}
              fill={
                entry.color ?? colors.charts.common.data[index % colors.charts.common.data.length]
              }
            />
          ))}
          {centerLabel && <Label content={<CenterLabel options={centerLabel} />} />}
        </Pie>
        {children && children}
        {chartOptions.legend && (
          <Legend layout="vertical" align="right" verticalAlign="top" iconType="plainline" />
        )}
        {chartOptions.tooltip && <Tooltip />}
      </RechartsPieChart>
    );

    return chartOptions.responsive ? (
      <ResponsiveContainer {...chartOptions.responsiveDimensions}>{chart}</ResponsiveContainer>
    ) : (
      chart
    );
  }
}

PieChart.contextType = ThemeContext;

export { PieChart };
