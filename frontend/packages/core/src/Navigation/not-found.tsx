import React from "react";
import ThumbDownIcon from "@mui/icons-material/ThumbDown";
import { Typography } from "@mui/material";

import { Grid } from "../Layout";
import { styled } from "../Utils";

const Container = styled(Grid)({
  minHeight: "80vh",
});

const IconContainer = styled(Grid)({
  color: "#02acbe",
  fontSize: "7rem",
});

const NotFound: React.FC<{}> = () => (
  <Container container direction="column" justifyContent="center" alignItems="center">
    <IconContainer item>
      <ThumbDownIcon fontSize="inherit" />
    </IconContainer>
    <Grid item>
      <Typography align="center" color="textPrimary" variant="h3">
        <Grid item>Whoops...</Grid>
        <Grid item>Looks like you took a wrong turn</Grid>
      </Typography>
      <Typography align="center" color="textPrimary" variant="h6">
        &lt; 404 Not Found &gt;
      </Typography>
    </Grid>
  </Container>
);

export default NotFound;
