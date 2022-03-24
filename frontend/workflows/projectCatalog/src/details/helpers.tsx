import React from "react";
import type { ReactTimeagoProps as TimeAgoProps } from "react-timeago";
import TimeAgo from "react-timeago";

const unitformatter = (unit: string): string => {
  switch (unit) {
    case "month":
      // month -> mo
      return unit.substring(0, 2);
    default:
      return unit.charAt(0);
  }
};

interface EventTimeProps extends Pick<TimeAgoProps, "date"> {
  onClick?: () => void;
}

const EventTime = ({ onClick, ...props }: EventTimeProps) => (
  <TimeAgo
    {...props}
    formatter={(value, unit) => `${value}${unitformatter(unit)}`}
    onClick={onClick}
  />
);

const parseTimestamp = (timestamp?: number | Long | null): number => {
  return parseInt(timestamp?.toString() || "0", 10);
};

const setMilliseconds = (timestamp?: number | Long | null): number => {
  const ts = new Date(0);
  return ts.setUTCMilliseconds(parseTimestamp(timestamp));
};

export { EventTime, unitformatter, setMilliseconds };
