import * as React from "react";
import type { CardHeaderSummaryProps, ClutchError } from "@clutch-sh/core";
import { Card as ClutchCard, CardContent, CardHeader, Error, IconButton } from "@clutch-sh/core";
import styled from "@emotion/styled";
import { Grid, LinearProgress } from "@material-ui/core";
import KeyboardArrowDownIcon from "@material-ui/icons/KeyboardArrowDown";
import KeyboardArrowUpIcon from "@material-ui/icons/KeyboardArrowUp";

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
  summary?: CardHeaderSummaryProps[];
  title?: React.ReactNode & string;
}

const Card = ({ avatar, children, error, isLoading, summary, title }: CardProps) => {
  const [expanded, setExpanded] = React.useState(true);

  const handleExpandClick = () => {
    setExpanded(!expanded);
  };

  return (
    <Grid item xs={12} sm={12} md={12} lg={6}>
      <ClutchCard>
        <CardHeader
          actions={
            <IconButton onClick={handleExpandClick} size="small" variant="neutral">
              {expanded ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
            </IconButton>
          }
          avatar={avatar}
          summary={summary}
          title={title}
        >
          <StyledProgressContainer>
            {isLoading && <LinearProgress color="secondary" />}
          </StyledProgressContainer>
        </CardHeader>
        {expanded && (
          <CardContent padding={0} collapsible maxHeight={500}>
            {error ? <Error subject={error} /> : children}
          </CardContent>
        )}
      </ClutchCard>
    </Grid>
  );
};

export default Card;
