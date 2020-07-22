import React from "react";
import { Grid, Typography } from "@material-ui/core";
import ThumbUpIcon from "@material-ui/icons/ThumbUp";
import styled from "styled-components";

const GridIcon = styled(Grid)`
  ${({ theme }) => `
  padding-top: 20px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  color: ${theme.palette.accent.main};
  font-size: 7rem;
  `}
`;

const Icon = styled(ThumbUpIcon)`
  font-size: 0.5em;
  margin-bottom: 10px;
`;

const Confirmation: React.FC<{ action: string }> = ({ action, children }) => (
  <Grid container direction="column" justify="center" alignItems="center">
    <GridIcon item>
      <Icon />
    </GridIcon>
    <Grid item>
      <Typography align="center" color="textPrimary" variant="h5">
        <Grid item>{action} requested!</Grid>
      </Typography>
      {children}
    </Grid>
  </Grid>
);

export default Confirmation;
