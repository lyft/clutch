import React, { PureComponent } from "react";
import {
  Cell,
  Legend,
  Pie,
  PieChart as RechartsPieChart,
  ResponsiveContainer,
  Sector,
  Tooltip,
} from "recharts";

import type { PieChartData } from "./types";

interface PieChartProps {
  data: PieChartData[];
  options?: {
    dimensions?: {
      height?: number;
      width?: number;
      innerRadius?: number;
      outerRadius?: number;
      paddingAngle?: number;
      cx?: string;
      cy?: string;
    };
    label?: boolean | React.ReactElement;
    labelLine?: boolean;
    legend?: boolean;
    responsive?: boolean;
    activeTooltip?: boolean;
    tooltip?: boolean;
  };
  children?: React.ReactChild;
}

interface PieChartState {
  activeIndex?: number;
}

const DEFAULT_COLORS = [
  "#3548D4",
  "#40A05A",
  "#B09027",
  "#D87313",
  "#C2302E",
  "#0D1030",
  "#8884D8",
];

const renderActiveShape = props => {
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
      <text x={cx} y={cy} dy={8} textAnchor="middle">
        {payload.activeLabel ?? payload.name}
      </text>
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
      <text
        x={ex + (cos >= 0 ? 1 : -1) * 12}
        y={ey}
        textAnchor={textAnchor}
        fill="#333"
      >{`${payload.name} ${value}`}</text>
      <text x={ex + (cos >= 0 ? 1 : -1) * 12} y={ey} dy={18} textAnchor={textAnchor} fill="#999">
        {`(${(percent * 100).toFixed(2)}%)`}
      </text>
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
    const { children, data, options } = this.props;

    const chartOptions = {
      activeTooltip: true,
      responsive: true,
      ...(options || {}),
      dimensions: {
        height: 275,
        width: 275,
        innerRadius: 60,
        outerRadius: 80,
        paddingAngle: 2,
        cx: "50%",
        cy: "50%",
        ...(options?.dimensions || {}),
      },
    };

    const additionalProps = {
      ...(chartOptions.activeTooltip
        ? // eslint-disable-next-line react/destructuring-assignment
          { activeIndex: this.state?.activeIndex, activeShape: renderActiveShape }
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
          fill="#8884d8"
          dataKey="value"
          onMouseEnter={this.onPieEnter}
          {...chartOptions.dimensions}
          {...additionalProps}
        >
          {data.map((entry, index) => (
            <Cell
              // eslint-disable-next-line react/no-array-index-key
              key={`cell-${index}`}
              fill={entry.color ?? DEFAULT_COLORS[index % DEFAULT_COLORS.length]}
            />
          ))}
        </Pie>
        {children && children}
        {chartOptions.legend && (
          <Legend
            layout="vertical"
            align="right"
            verticalAlign="middle"
            iconType="plainline"
            height={36}
          />
        )}
        {chartOptions.tooltip && <Tooltip />}
      </RechartsPieChart>
    );

    if (chartOptions.responsive) {
      return (
        <ResponsiveContainer width="99%" height={chartOptions.dimensions.height}>
          {chart}
        </ResponsiveContainer>
      );
    }

    return chart;
  }
}

export default PieChart;
