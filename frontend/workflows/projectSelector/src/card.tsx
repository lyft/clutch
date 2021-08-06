import * as React from "react";
import type { ClutchError } from "@clutch-sh/core";
import { Card as ClutchCard, CardHeader, Error } from "@clutch-sh/core";
import styled from "@emotion/styled";
import { Grid, LinearProgress } from "@material-ui/core";

const StyledProgressContainer = styled.div({
  height: "4px",
  ".MuiLinearProgress-root": {
    backgroundColor: "rgb(194, 200, 242)",
  },
  ".MuiLinearProgress-bar": {
    backgroundColor: "#3548D4",
  },
});

interface CardProps {
  avatar?: React.ReactNode;
  children: React.ReactNode;
  error?: ClutchError;
  isLoading?: boolean;
  title?: React.ReactNode & string;
}

const Card = ({ avatar, children, error, isLoading, title }: CardProps) => (
  <Grid item xs={12} sm={6}>
    <ClutchCard>
      <CardHeader avatar={avatar} title={title} />
      {isLoading && (
        <StyledProgressContainer>
          <LinearProgress color="secondary" />
        </StyledProgressContainer>
      )}
      {error ? <Error subject={error} /> : children}
    </ClutchCard>
  </Grid>
);

export default Card;
