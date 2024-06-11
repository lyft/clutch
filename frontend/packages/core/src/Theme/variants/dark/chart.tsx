import type { ChartColors } from "../../types";

const COLORS: string[] = [
  "#651FFF",
  "#FF4081",
  "#0091EA",
  "#00695C",
  "#9E9D24",
  "#880E4F",
  "#01579B",
  "#F4511E",
  "#009688",
  "#C2185B",
  "#1A237E",
  "#7C4DFF",
  "#88451D",
  "#AA00FF",
];

const chartColors: ChartColors = {
  common: {
    data: COLORS,
  },
  pie: {
    labelPrimary: "#0D1030",
    labelSecondary: "#8D8E9E",
  },
  linearTimeline: {
    xAxisStroke: "#FFF",
    tooltipBackgroundColor: "#FFF",
    tooltipTextColor: "#000",
    gridBackgroundColor: "#FFF",
    gridStroke: "#000",
  },
};

export default chartColors;
