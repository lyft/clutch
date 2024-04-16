import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { client, Grid, Link, styled, TimeAgo as EventTime, Typography } from "@clutch-sh/core";
import { faClock } from "@fortawesome/free-regular-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import * as _ from "lodash";

const StyledLink = styled(Link)({
  whiteSpace: "nowrap",
});

const fetchProjectInfo = (
  project: string,
  allowDisabled: boolean = false
): Promise<IClutch.core.project.v1.IProject> =>
  client
    .post("/v1/project/getProjects", { projects: [project], excludeDependencies: true })
    .then(resp => {
      const { results = {}, partialFailures } = resp.data as IClutch.project.v1.GetProjectsResponse;

      const projectResults = _.get(results, [project, "project"], {});

      if (_.isEmpty(projectResults) && allowDisabled && partialFailures) {
        // Will add disabled projects to the list of projects to display if requested
        const failuresMap = partialFailures
          .map(p => {
            if ((_.get(p, "message", "") ?? "").includes("disabled")) {
              const disabledProject = _.get(p, "details.[0]");
              _.set(
                disabledProject,
                ["data", "description"],
                `${disabledProject.name} is disabled.`
              );
              return disabledProject;
            }
            return null;
          })
          .filter(Boolean);

        if (failuresMap.length) {
          return failuresMap[0];
        }
      }

      return projectResults;
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
