import React from "react";
import { Grid, Typography } from "@material-ui/core";
import ThumbDownIcon from "@material-ui/icons/ThumbDown";
import styled from "styled-components";

const GridIcon = styled(Grid)`
  ${({ theme }) => `
  padding-top: 25%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  color: ${theme.palette.accent.main};
  font-size: 7rem;
  `}
`;

const NotFound: React.FC<{}> = () => (
  <Grid container direction="column" justify="center" alignItems="center">
    <GridIcon item>
      <ThumbDownIcon fontSize="inherit" />
    </GridIcon>
    <Grid item>
      <Typography align="center" color="textPrimary" variant="h3">
        <Grid item>Whoops...</Grid>
        <Grid item>Looks like you took a wrong turn</Grid>
      </Typography>
      <Typography align="center" color="textPrimary" variant="h6">
        &lt; 404 Not Found &gt;
      </Typography>
    </Grid>
  </Grid>
);

export default NotFound;
