import React from "react";
import { render, waitFor } from "@testing-library/react";

import "@testing-library/jest-dom";

import { TimeseriesChart } from "..";

jest.mock("recharts", () => {
  const OriginalModule = jest.requireActual("recharts");
  return {
    ...OriginalModule,
    ResponsiveContainer: ({ children }) => (
      <OriginalModule.ResponsiveContainer width={400} height={200}>
        {children}
      </OriginalModule.ResponsiveContainer>
    ),
  };
});

const mockDataSingleLine = [
  {
    lineName: "line1",
    timestamp: 1546300800000,
    value: 5,
  },
  {
    lineName: "line1",
    timestamp: 1546300900000,
    value: 20,
  },
  {
    lineName: "line1",
    timestamp: 1546301000000,
    value: 30,
  },
  {
    lineName: "line1",
    timestamp: 1546300700000,
    value: 5,
  },
  {
    lineName: "line1",
    timestamp: 1546300600000,
    value: 20,
  },
  {
    lineName: "line1",
    timestamp: 1546301800000,
    value: 30,
  },
];

const mockLines = [
  {
    dataKey: "value",
    color: "red",
    animationDuration: 0,
  },
];

const setup = () =>
  render(
    <TimeseriesChart
      data={mockDataSingleLine}
      xAxisDataKey="timestamp"
      yAxisDataKey="value"
      lines={mockLines}
      xDomainSpread={0.3}
      yDomainSpread={0.3}
      regularIntervalTicks
    />
  );

test("renders correctly", () => {
  const { container } = setup();

  expect(container.querySelector(".recharts-responsive-container")).toBeVisible();
});

test("renders a Cartesian Grid", async () => {
  const { container } = setup();

  await waitFor(() => {
    expect(container.querySelector(".recharts-cartesian-grid")).toBeDefined();
  });
});

test("renders a horizontal Cartesian Grid", async () => {
  const { container } = setup();

  await waitFor(() => {
    expect(container.querySelector(".recharts-cartesian-grid-horizontal")).toBeDefined();
    expect(
      container.querySelectorAll(".recharts-cartesian-grid-horizontal line").length
    ).toBeGreaterThan(0);
  });
});

test("renders a vertical Cartesian Grid", async () => {
  const { container } = setup();

  await waitFor(() => {
    expect(container.querySelector(".recharts-cartesian-grid-vertical")).toBeDefined();
    expect(
      container.querySelectorAll(".recharts-cartesian-grid-vertical line").length
    ).toBeGreaterThan(0);
  });
});

test("renders xAxis with type number property", async () => {
  const { container } = setup();

  await waitFor(() => {
    expect(
      container.querySelector(".recharts-xAxis .recharts-cartesian-axis-line")
    ).toHaveAttribute("type", "number");
  });
});

test("renders yAxis with type number property", async () => {
  const { container } = setup();

  await waitFor(() => {
    expect(
      container.querySelector(".recharts-yAxis .recharts-cartesian-axis-line")
    ).toHaveAttribute("type", "number");
  });
});

test("renders the x axis with 5 ticks", async () => {
  const { container } = setup();

  await waitFor(() => {
    expect(
      container.querySelector(".recharts-xAxis .recharts-cartesian-axis-ticks")?.childElementCount
    ).toBe(5);
  });
});

test("renders the y axis with 5 ticks", async () => {
  const { container } = setup();

  await waitFor(() => {
    expect(
      container.querySelector(".recharts-yAxis .recharts-cartesian-axis-ticks")?.childElementCount
    ).toBe(5);
  });
});

test("renders one line", async () => {
  const { container } = setup();

  await waitFor(() => {
    const chartLines = container.querySelectorAll(".recharts-line-curve");
    expect(chartLines).toHaveLength(1);
  });
});

test("renders a line with line type linear", async () => {
  const { container } = setup();

  let chartLine;
  await waitFor(() => {
    [chartLine] = container.querySelectorAll(".recharts-line-curve");
  });

  expect(chartLine).toHaveAttribute("type", "linear");
});

test("renders a line with red stroke", async () => {
  const { container } = setup();

  let chartLine;
  await waitFor(() => {
    [chartLine] = container.querySelectorAll(".recharts-line-curve");
  });

  expect(chartLine).toHaveAttribute("stroke", "red");
});
