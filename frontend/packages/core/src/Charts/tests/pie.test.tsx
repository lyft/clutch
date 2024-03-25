import React from "react";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import { ThemeProvider } from "../../Theme";
import { PieChart } from "../pie";
import type { PieChartData } from "../types";

const mockData: PieChartData[] = [
  { name: "Test Value #1", value: 1 },
  { name: "Test Value #2", value: 2 },
  { name: "Test Value #3", value: 3 },
  { name: "Test Value #4", value: 4 },
  { name: "Test Value #5", value: 5 },
  { name: "Test Value #6", value: 6 },
  { name: "Test Value #7", value: 7 },
];

jest.mock("recharts", () => {
  const OriginalModule = jest.requireActual("recharts");
  return {
    ...OriginalModule,
    ResponsiveContainer: ({ children }) => (
      <OriginalModule.ResponsiveContainer width={550} height={550}>
        {children}
      </OriginalModule.ResponsiveContainer>
    ),
  };
});

const setup = props =>
  render(
    <ThemeProvider>
      <PieChart data={mockData} {...props} />
    </ThemeProvider>
  );

test("renders correctly", () => {
  const { container } = setup({});

  expect(container.querySelector(".recharts-responsive-container")).toBeVisible();
});

test("renders correct number of pie sectors", () => {
  const { container } = setup({});

  expect(container.querySelectorAll(".recharts-pie-sector")).toHaveLength(mockData.length);
});

test("renders a legend", () => {
  const { container } = setup({ legend: true });

  expect(container.querySelector(".recharts-legend-wrapper")).toBeVisible();
  expect(container.querySelectorAll(".recharts-legend-item")).toHaveLength(mockData.length);
});
