import * as React from "react";
import type { ClutchError } from "@clutch-sh/core";
import { Card as ClutchCard, CardHeader, Error } from "@clutch-sh/core";
import { Grid } from "@material-ui/core";

interface CardProps {
  avatar?: React.ReactNode;
  title?: React.ReactNode & string;
  error?: ClutchError;
  children: React.ReactNode;
}

const Card = ({ avatar, title, error, children }: CardProps) => (
  <Grid item xs={12} sm={6}>
    <ClutchCard>
      <CardHeader avatar={avatar} title={title} />
      {error ? <Error subject={error} /> : children}
    </ClutchCard>
  </Grid>
);

export default Card;
