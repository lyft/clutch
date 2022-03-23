import React from "react";
import { Link, styled, TableRow, Typography } from "@clutch-sh/core";
import { Avatar } from "@material-ui/core";
import { capitalize } from "lodash";

import type { Incident, Urgency } from "./types";

const AVATAR_COLOR_PALETTE = ["#727FE1", "#32A140", "#F59E0B", "#E95F52"];

const StyledSpacer = styled("span")({
  height: "32px",
  width: "32px",
  marginRight: "8px",
});

const StyledAvatar = styled(Avatar)<{ color: string }>(
  {
    height: "32px",
    width: "32px",
    marginRight: "8px",
  },
  props => ({
    backgroundColor: props.color,
  })
);

const StyledAssignee = styled("span")({
  display: "flex",
  alignItems: "center",
  whiteSpace: "nowrap",
});

const StyledStatus = styled("span")<{ urgency: Urgency }>(
  {
    height: "7px",
    width: "7px",
    borderRadius: "50%",
    display: "inline-block",
    marginLeft: "60px",
  },
  props => ({
    backgroundColor: props.urgency === "HIGH" ? "#C2302E" : "#D87313",
  })
);

const IncidentRow = ({ incident }: { incident: Incident }) => {
  const oneDay = 24 * 60 * 60 * 1000;
  const today = new Date();
  const createdDate = new Date(incident.created);
  const dayDelta = Math.round((today.getTime() - createdDate.getTime()) / oneDay);

  const randomColor = AVATAR_COLOR_PALETTE[Math.floor(Math.random() * AVATAR_COLOR_PALETTE.length)];

  let userInitials = "";
  let assignee = "";
  if (incident.assignments.length > 0) {
    assignee = incident.assignments[0].assignee;
    const parts = assignee.split(" ");
    parts.forEach(part => {
      userInitials += part.charAt(0);
    });
  }

  const Description = () => (
    <Typography color="inherit" variant="body4">
      {incident.description}
    </Typography>
  );
  return (
    <TableRow key={incident.id}>
      <StyledStatus urgency={incident.urgency} />
      {incident.url ? (
        <Link href={incident.url}>
          <Description />
        </Link>
      ) : (
        <Description />
      )}
      <Typography variant="body4">Urgency:&nbsp;{capitalize(incident.urgency)}</Typography>
      <span style={{ whiteSpace: "nowrap" }}>
        <Typography variant="body4">({dayDelta} days ago)</Typography>
      </span>
      <StyledAssignee>
        {userInitials ? (
          <StyledAvatar color={randomColor}>
            <Typography variant="subtitle3" color="#ffffff">
              {userInitials}
            </Typography>
          </StyledAvatar>
        ) : (
          <StyledSpacer />
        )}
        <Typography variant="body3">{assignee || "No Assignee"}</Typography>
      </StyledAssignee>
    </TableRow>
  );
};

export default IncidentRow;
