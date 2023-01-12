import React from "react";
import type { ReactTimeagoProps as TimeAgoProps } from "react-timeago";
import ReactTimeago from "react-timeago";

interface EventTimeProps extends Pick<TimeAgoProps, "date"> {
  onClick?: () => void;
  short?: boolean;
  live?: boolean;
  formatter?: (value, unit, suffix) => string;
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
 * Will take a date/timestamp in and calculate the timeago for it
 * @param date Date is a date in the past or the future. This can be a Date Object, A UTC date-string or number of milliseconds since epoch time.
 * @param live (default) TimeAgo is live by default and will auto update it's value. However, if you don't want this behaviour, you can set live:false.
 * @param short (default) will shorten the unit string (day -> d)
 * @returns react component representing the timeago
 */
const TimeAgo = ({ short = true, onClick, formatter, ...props }: EventTimeProps) => (
  <ReactTimeago
    {...props}
    formatter={
      formatter ??
      ((value, unit) =>
        `${setMilliseconds(value)}${
          short ? unitFormatter(unit) : value > 1 ? ` ${unit}s` : ` ${unit}`
        }`)
    }
    onClick={onClick}
  />
);

export default TimeAgo;
