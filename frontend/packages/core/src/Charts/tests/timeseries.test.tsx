import React from "react";
import { shallow } from "enzyme";

import { TimeseriesChart } from "..";

describe("<TimeseriesChart />", () => {
  describe("basic rendering", () => {
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

    let timeseriesChart;
    beforeEach(() => {
      timeseriesChart = shallow(
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
    });

    it("renders", () => {
      expect(timeseriesChart.find(TimeseriesChart)).toBeDefined();
    });

    it("renders XAxis with type number property", () => {
      expect(timeseriesChart.find("XAxis").props().type).toBe("number");
    });

    it("renders YAxis with type number property", () => {
      expect(timeseriesChart.find("YAxis").props().type).toBe("number");
    });

    it("renders a Cartesian Grid", () => {
      expect(timeseriesChart.find("CartesianGrid")).toBeDefined();
    });

    it("renders a line with line type linear", () => {
      expect(timeseriesChart.find("Line").prop("type")).toBe("linear");
    });

    it("renders a line with red stroke", () => {
      expect(timeseriesChart.find("Line").prop("stroke")).toBe("red");
    });

    it("renders the x axis with 5 ticks", () => {
      expect(timeseriesChart.find("XAxis").prop("tickCount")).toBe(5);
    });

    it("renders the y axis with 5 ticks", () => {
      expect(timeseriesChart.find("YAxis").prop("tickCount")).toBe(5);
    });

    it("renders one line", () => {
      expect(timeseriesChart.find("Line")).toHaveLength(1);
    });
  });
});
