export const localTimeFormatter = (timestamp: number) => {
  return new Date(timestamp).toLocaleTimeString();
};

export const isoTimeFormatter = (timestamp: number) => {
  return new Date(timestamp).toISOString();
};

export const dateTimeFormatter = (timestamp: number) => {
  return new Date(timestamp).toDateString();
};

export const getMinAndMaxOfRangeUsingKey = (data: any, key: string) => {
  if (!data) {
    return {
      min: null,
      max: null,
    };
  }
  const filtered = data.filter(d => d[key]).map(d => d[key]);
  return { min: Math.min(...filtered), max: Math.max(...filtered) };
};

// Edge ratio refers to the multiplicative part of the amount that will be added to the max,
// and subtracted from the min.
export const calculateDomainEdges = (data: any, dataKey: string, edgeRatio: number) => {
  // Get the max and min of the domain, then calculate a certain amount`out from each edge.
  const { min, max } = getMinAndMaxOfRangeUsingKey(data, dataKey);
  if (
    edgeRatio <= 0 ||
    min === null ||
    max === null ||
    min === Infinity ||
    max === -Infinity ||
    min === undefined ||
    max === undefined ||
    Number.isNaN(min) ||
    Number.isNaN(max)
  ) {
    return [min, max];
  }
  if (max === min) {
    const minEdge = min - min * edgeRatio;
    const maxEdge = max + max * edgeRatio;
    return [minEdge, maxEdge];
  }
  const edgeVal = (max - min) * edgeRatio;
  const minEdge = min - edgeVal;
  const maxEdge = max + edgeVal;
  return [minEdge, maxEdge];
};

const oneSec = 1000; // in ms
const fifteenSeconds = 15 * oneSec;
const oneMin = 60 * oneSec;
const threeMins = 3 * oneMin;
const fiveMins = 5 * oneMin;
const tenMins = 10 * oneMin;
const fifteenMins = 15 * oneMin;
const halfHour = 30 * oneMin;
const oneHour = 60 * oneMin;
const threeHours = 3 * oneHour;
const sixHours = 6 * oneHour;
const twelveHours = 12 * oneHour;
const dayDuration = 24 * oneHour;
const weekDuration = 7 * dayDuration;
const monthDuration = 30 * dayDuration;
const yearDuration = 365 * dayDuration;

const zoomLevelsToIntervals = {
  oneMin: fifteenSeconds,
  threeMins: fifteenSeconds,
  fiveMins: oneMin,
  tenMins: oneMin,
  fifteenMins: threeMins,
  halfHour: fiveMins,
  oneHour: tenMins,
  threeHours: halfHour,
  sixHours: oneHour,
  twelveHours: threeHours,
  day: sixHours,
  week: dayDuration,
  month: weekDuration,
  year: monthDuration,
};

// This function allows us to get the starting point for our ticks, as well as the space between ticks.
// We have presets according to the span between the min and max timestamps.
export const getLeftSideAndIntervalForTicks = (min: number, max: number) => {
  if (
    min === null ||
    max === null ||
    min === Infinity ||
    max === -Infinity ||
    min === undefined ||
    max === undefined ||
    Number.isNaN(min) ||
    Number.isNaN(max) ||
    min > max
  ) {
    return {
      leftSide: null,
      interval: null,
    };
  }

  const diff = max - min;
  let zoomLevel = "";
  switch (true) {
    case diff < oneMin:
      zoomLevel = "oneMin";
      break;
    case diff < fiveMins:
      zoomLevel = "fiveMins";
      break;
    case diff < tenMins:
      zoomLevel = "tenMins";
      break;
    case diff < fifteenMins:
      zoomLevel = "fifteenMins";
      break;
    case diff < halfHour:
      zoomLevel = "halfHour";
      break;
    case diff < oneHour:
      zoomLevel = "oneHour";
      break;
    case diff < threeHours:
      zoomLevel = "threeHours";
      break;
    case diff < sixHours:
      zoomLevel = "sixHours";
      break;
    case diff < twelveHours:
      zoomLevel = "twelveHours";
      break;
    case diff < dayDuration:
      zoomLevel = "day";
      break;
    case diff < weekDuration:
      zoomLevel = "week";
      break;
    case diff < monthDuration:
      zoomLevel = "month";
      break;
    case diff < yearDuration:
      zoomLevel = "month";
      break;
    default:
      zoomLevel = "year";
  }

  const interval = zoomLevelsToIntervals[zoomLevel];
  const leftSide = min - (min % interval);
  return { leftSide, interval };
};

// Based off the min and max, calculate where the regularly spaced tick marks should be.
// We modulo down to the closest similar timestamp based off the difference between the min and max.
// We then iterate from that value to the max pushing timestamps to our array when we land on an
// appropriate interval.
export const calculateTicks = (data: any, dataKey: string) => {
  const { min, max } = getMinAndMaxOfRangeUsingKey(data, dataKey);
  const { leftSide, interval } = getLeftSideAndIntervalForTicks(min, max);
  const ticks = [];

  if (
    leftSide === null ||
    interval === null ||
    min === null ||
    max === null ||
    min === Infinity ||
    max === -Infinity ||
    min === undefined ||
    max === undefined ||
    Number.isNaN(min) ||
    Number.isNaN(max)
  ) {
    return ticks;
  }

  for (let i = leftSide; i <= max; i += interval) {
    ticks.push(i);
  }

  return ticks;
};

export interface CustomTooltipProps {
  active: boolean;
  payload: any; // A huge object that contains all the info for the data point and more
}
