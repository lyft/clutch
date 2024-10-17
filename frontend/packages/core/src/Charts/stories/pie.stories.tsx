import React from "react";
import type { Meta } from "@storybook/react";

import styled from "../../styled";
import type { PieChartProps } from "../pie";
import { PieChart as PieChartComponent } from "..";

export default {
  title: "Core/Charts/Pie",
  component: PieChartComponent,
} as Meta;

const ChartContainer = styled("div")({
  height: "550px",
  width: "550px",
});

const Template = (props: PieChartProps) => (
  <ChartContainer>
    <PieChartComponent {...props} />
  </ChartContainer>
);

export const PieChart = Template.bind({});
PieChart.args = {
  data: [
    { name: "Test Value #1", value: 1 },
    { name: "Test Value #2", value: 2 },
    { name: "Test Value #3", value: 3 },
    { name: "Test Value #4", value: 4 },
    { name: "Test Value #5", value: 5 },
    { name: "Test Value #6", value: 6 },
    { name: "Test Value #7", value: 7 },
  ],
};
