import * as React from "react";
import type { CardHeaderSummaryProps, ClutchError } from "@clutch-sh/core";
import {
  Card as ClutchCard,
  CardContent,
  CardHeader,
  Error,
  Grid,
  IconButton,
  styled,
} from "@clutch-sh/core";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import KeyboardArrowUpIcon from "@mui/icons-material/KeyboardArrowUp";
import { LinearProgress, Theme } from "@mui/material";

const StyledProgressContainer = styled("div")(({ theme }: { theme: Theme }) => ({
  height: "4px",
  ".MuiLinearProgress-root": {
    backgroundColor: theme.palette.primary[400],
  },
  ".MuiLinearProgress-bar": {
    backgroundColor: theme.palette.primary[600],
  },
}));

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
