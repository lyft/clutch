import React from "react";
import { Card, Grid, styled } from "@clutch-sh/core";
import { faGithub, faSlack } from "@fortawesome/free-brands-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

const StyledContainer = styled(Grid)({
  padding: "20px",
});

const Details = () => (
  <StyledContainer container direction="row" wrap="nowrap">
    {/* Column for project details and header */}
    <Grid container item direction="column">
      <Grid item style={{ marginBottom: "22px" }} />
    </Grid>
    {/* Column for project quick links */}
    <Grid container item direction="column" xs={1}>
      <Card>
        <Grid
          container
          item
          direction="column"
          alignItems="center"
          spacing={1}
          style={{ padding: "7px 0" }}
        >
          <Grid item>
            <FontAwesomeIcon icon={faGithub} size="3x" />
          </Grid>
          <Grid item>
            <FontAwesomeIcon icon={faSlack} size="3x" />
          </Grid>
        </Grid>
      </Card>
    </Grid>
  </StyledContainer>
);

export default Details;
