import React from "react";
import { Link } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Card, Grid } from "@clutch-sh/core";

const QuickLinksCard = (linkGroups: IClutch.core.project.v1.ILinkGroup[]) => {
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
        {linkGroups?.map(linkGroup => {
          return linkGroup.links?.map(link => {
            return (
              <Grid item key={link.name}>
                <Link to={link.url}>
                  <img width="29px" height="29px" src={linkGroup.imagePath} alt={link.name} />
                </Link>
              </Grid>
            );
          });
        })}
      </Grid>
    </Card>
  );
};

export default QuickLinksCard;
