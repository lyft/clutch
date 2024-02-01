import React from "react";
import styled from "@emotion/styled";
import ThumbDownIcon from "@mui/icons-material/ThumbDown";
import { Grid, Theme, Typography } from "@mui/material";

const Container = styled(Grid)`
  minheight: 80vh;
`;

const IconContainer = styled(Grid)(({ theme }: { theme: Theme }) => ({
  color: theme.palette.brandColor,
  fontSize: "7rem",
}));

const NotFound: React.FC<{}> = () => (
  <Container
    container
    direction="column"
    justifyContent="center"
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
