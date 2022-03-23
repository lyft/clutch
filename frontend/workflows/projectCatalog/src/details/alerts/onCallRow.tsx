import React from "react";
import { Grid, styled, Typography } from "@clutch-sh/core";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Avatar } from "@material-ui/core";

import { LinkText, StyledLink } from "../card";

import type { OnCall, User } from "./types";

const StyledAvatar = styled(Avatar)({
  width: "20px",
  height: "20px",
  fontSize: "8px",
  alignSelf: "center",
  backgroundColor: "#727FE1",
});

const OnCallUser = ({ name, url }: User) => {
  const initials = name.split(" ").map(n => n.charAt(0));
  const userName = (
    <Typography variant="body3" color="#3548D4">
      {name}
    </Typography>
  );

  return (
    <Grid item>
      <Grid container item justify="flex-start" alignItems="center" spacing={1}>
        <Grid item>
          <StyledAvatar>{initials}</StyledAvatar>
        </Grid>
        <Grid item>{url ? <StyledLink href={url}>{userName}</StyledLink> : userName}</Grid>
      </Grid>
    </Grid>
  );
};

const OnCallRow = ({ text, icon, users, url }: OnCall) => (
  <>
    {text && url && (
      <Grid container spacing={1}>
        {icon && (
          <Grid item style={{ paddingLeft: "10px" }}>
            <FontAwesomeIcon icon={icon} size="lg" />
          </Grid>
        )}
        <Grid item>
          <LinkText text={text} link={url} />
        </Grid>
      </Grid>
    )}
    {users?.length && (
      <Grid item>
        <Grid container direction="row" justify="flex-start" alignItems="center" spacing={1}>
          {users.map((user: User) => (
            <OnCallUser {...user} />
          ))}
        </Grid>
      </Grid>
    )}
  </>
);

export default OnCallRow;
