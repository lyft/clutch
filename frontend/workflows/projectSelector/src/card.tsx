import * as React from "react";
import type { CardHeaderSections, ClutchError } from "@clutch-sh/core";
import { Card as ClutchCard, CardContent, CardHeader, Error } from "@clutch-sh/core";
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
  sections?: CardHeaderSections[];
  title?: React.ReactNode & string;
}

const Card = ({ avatar, children, error, isLoading, sections, title }: CardProps) => (
  <Grid item xs={12} sm={6}>
    <ClutchCard>
      <CardHeader avatar={avatar} sections={sections} title={title}>
        <StyledProgressContainer>
          {isLoading && <LinearProgress color="secondary" />}
        </StyledProgressContainer>
      </CardHeader>
      <CardContent padding={0}>{error ? <Error subject={error} /> : children}</CardContent>
    </ClutchCard>
  </Grid>
);

export default Card;
