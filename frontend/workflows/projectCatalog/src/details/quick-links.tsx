import React from "react";
import { Card, Grid } from "@clutch-sh/core";
import { faGithub, faSlack } from "@fortawesome/free-brands-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

const QuickLinksCard = () => {
  return (
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
  );
};

export default QuickLinksCard;
