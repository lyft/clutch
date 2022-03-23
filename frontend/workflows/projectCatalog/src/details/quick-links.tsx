import React from "react";
import { Card, Grid, styled } from "@clutch-sh/core";
import { faGithub, faSlack } from "@fortawesome/free-brands-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

const StyledCard = styled(Card)({
  width: "100%",
  height: "fit-content",
  padding: "15px",
});

const QuickLinksCard = () => {
  return (
    <StyledCard container direction="column" justify="space-evenly" alignItems="center">
      <Grid item>
        <FontAwesomeIcon icon={faGithub} size="3x" />
      </Grid>
      <Grid item>
        <FontAwesomeIcon icon={faSlack} size="3x" />
      </Grid>
    </StyledCard>
  );
};

export default QuickLinksCard;
