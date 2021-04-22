import React from "react";
import styled from "@emotion/styled";
import { Grid, Typography } from "@material-ui/core";
import ThumbDownIcon from "@material-ui/icons/ThumbDown";

const Container = styled(Grid)`
  minheight: 80vh;
`;

const IconContainer = styled(Grid)({
  color: "#02acbe",
  fontSize: "7rem",
});

const NotFound: React.FC<{}> = () => (
  <Container
    container
    direction="column"
    justify="center"
    alignItems="center"
    style={{ minHeight: "80vh" }}
  >
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
