import * as React from "react";
import { Grid } from "@material-ui/core";

import Error from "../Feedback/error/index";
import type { ClutchError } from "../Network/errors";

import { Card as ClutchCard, CardHeader } from "./card";

interface WorkflowCardProps {
  avatar?: React.ReactNode;
  title?: React.ReactNode & string;
  error?: ClutchError;
  children: React.ReactNode;
}

const WorkflowCard = ({ avatar, title, error, children }: WorkflowCardProps) => (
  <Grid item xs={12} sm={7}>
    <ClutchCard>
      <CardHeader avatar={avatar} title={title} />
      {error ? <Error subject={error} /> : children}
    </ClutchCard>
  </Grid>
);

export default WorkflowCard;
