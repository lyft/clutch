import React from "react";
import type { ReactTimeagoProps as TimeAgoProps } from "react-timeago";
import TimeAgo from "react-timeago";
import { Typography } from "@clutch-sh/core";

const getRepoName = (project: string, metadata: any): string => {
  const repoGitUrl = metadata[project]?.data?.repository;
  if (!repoGitUrl) {
    // we don't have a repo git url link, so we default to using project name
    return project;
  }

  // can test at https://regex101.com/r/YkZUiA/1
  const repoGitUrlRegex = /(\w+)\/([^.]+)/g;
  const repoName = repoGitUrl.match(repoGitUrlRegex)?.[0];
  if (!repoName) {
    // we didn't get a match, so we default to using project name
    return project;
  }

  return repoName;
};

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

// displayed in the card header summary when there aren't selected projects
const DefaultSummaryTitle = () => <Typography variant="subtitle2">-</Typography>;

export { DefaultSummaryTitle, getRepoName, EventTime, unitformatter, setMilliseconds };
