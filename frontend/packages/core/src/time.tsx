import type { ReactTimeagoProps as TimeAgoProps } from "react-timeago";
import React from "react";
import TimeAgo from "react-timeago";

interface EventTimeProps extends Pick<TimeAgoProps, "date"> {
  onClick?: () => void;
  short?: boolean;
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

const EventTime = ({ short = true, onClick, ...props }: EventTimeProps) => (
  <TimeAgo
    {...props}
    formatter={(value, unit) =>
      `${setMilliseconds(value)}${
        short ? unitFormatter(unit) : value > 1 ? ` ${unit}s` : ` ${unit}`
      }`
    }
    onClick={onClick}
  />
);

export { EventTime };
