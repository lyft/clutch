import React from "react";
import { Grid } from "@clutch-sh/core";
import { faGithub, faSlack } from "@fortawesome/free-brands-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

import { StyledCard } from "./cards/base";

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
