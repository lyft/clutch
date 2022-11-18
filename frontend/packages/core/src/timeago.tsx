import React from "react";
import type { ReactTimeagoProps as TimeAgoProps } from "react-timeago";
import ReactTimeago from "react-timeago";

interface EventTimeProps extends Pick<TimeAgoProps, "date"> {
  onClick?: () => void;
  short?: boolean;
  live?: boolean;
}

const unitFormatter = (unit: string): string => {
  switch (unit) {
    case "month":
      // month -> mo
      return unit.substring(0, 2);
    default:
      return unit.charAt(0);
  }
};

const parseTimestamp = (timestamp?: number | Long | null): number => {
  return parseInt(timestamp?.toString() || "0", 10);
};

const setMilliseconds = (timestamp?: number | Long | null): number => {
  const ts = new Date(0);
  return ts.setUTCMilliseconds(parseTimestamp(timestamp));
};

/**
 * Will take a millsecond timestamp in and calculate the timeago for it
 * @param date the millisecond timestamp
 * @param short (default) will shorten the unit string (day -> d)
 * @param live (default) will auto increment based on a given time
 * @returns react component representing the timeago
 */
const TimeAgo = ({ short = true, onClick, ...props }: EventTimeProps) => (
  <ReactTimeago
    {...props}
    formatter={(value, unit) =>
      `${setMilliseconds(value)}${
        short ? unitFormatter(unit) : value > 1 ? ` ${unit}s` : ` ${unit}`
      }`
    }
    onClick={onClick}
  />
);

export default TimeAgo;
