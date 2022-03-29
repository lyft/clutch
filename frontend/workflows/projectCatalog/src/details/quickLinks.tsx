import React from "react";
import { Link } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Card, Grid } from "@clutch-sh/core";

const QuickLinks = (project: IClutch.core.project.v1.IProject) => {
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
        {project?.linkGroups?.map(linkGroup => {
          return linkGroup.links?.map(link => {
            return (
              <Grid item key={link.name}>
                <Link
                  to={link.url}
                  style={{
                    textDecoration: "none",
                    color: "inherit",
                  }}
                >
                  <img src={linkGroup.imagePath} alt={link.name} />
                </Link>
              </Grid>
            );
          });
        })}
      </Grid>
    </Card>
  );
};

export default QuickLinks;
