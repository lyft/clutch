import React from "react";
import type { ReactTimeagoProps as TimeAgoProps } from "react-timeago";
import TimeAgo from "react-timeago";
import type { clutch as IClutch } from "@clutch-sh/api";
import { client, Grid, Link, styled, Typography } from "@clutch-sh/core";
import { faClock } from "@fortawesome/free-regular-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

const StyledLink = styled(Link)({
  whiteSpace: "nowrap",
});

const fetchProjectInfo = (project: string): Promise<IClutch.core.project.v1.IProject> =>
  client
    .post("/v1/project/getProjects", { projects: [project], excludeDependencies: true })
    .then(resp => {
      const { results = {} } = resp.data as IClutch.project.v1.GetProjectsResponse;

      return results[project] ? results[project].project ?? {} : {};
    });

const unitFormatter = (unit: string): string => {
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
    formatter={(value, unit) => `${value}${unitFormatter(unit)}`}
    onClick={onClick}
  />
);

const LinkText = ({ text, link }: { text: string; link?: string }) => {
  const returnText = <Typography variant="body2">{text}</Typography>;

  if (link && text) {
    return <StyledLink href={link}>{returnText}</StyledLink>;
  }

  return returnText;
};

const parseTimestamp = (timestamp?: number | Long | null): number => {
  return parseInt(timestamp?.toString() || "0", 10);
};

const setMilliseconds = (timestamp?: number | Long | null): number => {
  const ts = new Date(0);
  return ts.setUTCMilliseconds(parseTimestamp(timestamp));
};

const LastEvent = ({ time, ...props }: { time: number }) => {
  return time ? (
    <>
      <Grid item>
        <FontAwesomeIcon icon={faClock} />
      </Grid>
      <Grid item>
        <Typography variant="body4">
          <EventTime date={setMilliseconds(time)} {...props} /> ago
        </Typography>
      </Grid>
    </>
  ) : null;
};

export { fetchProjectInfo, LastEvent, LinkText };
