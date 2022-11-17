import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { client, EventTime, Grid, Link, styled, Typography } from "@clutch-sh/core";
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

const LinkText = ({ text, link }: { text: string; link?: string }) => {
  const returnText = <Typography variant="body2">{text}</Typography>;

  if (link && text) {
    return <StyledLink href={link}>{returnText}</StyledLink>;
  }

  return returnText;
};

const LastEvent = ({ time, ...props }: { time: number }) => {
  return time ? (
    <>
      <Grid item>
        <FontAwesomeIcon icon={faClock} />
      </Grid>
      <Grid item>
        <Typography variant="body4">
          <EventTime date={time} {...props} /> ago
        </Typography>
      </Grid>
    </>
  ) : null;
};

export { fetchProjectInfo, LastEvent, LinkText };
