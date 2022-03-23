import React from "react";
import { Card, ClutchError, Error, Grid, Link, styled, Typography } from "@clutch-sh/core";
import { faClock } from "@fortawesome/free-regular-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { LinearProgress } from "@material-ui/core";

import { EventTime, setMilliseconds } from "./helpers";

const StyledCard = styled(Card)({
  width: "100%",
  height: "fit-content",
  padding: "15px",
});

const StyledProgressContainer = styled("div")({
  height: "4px",
  ".MuiLinearProgress-root": {
    backgroundColor: "rgb(194, 200, 242)",
  },
  ".MuiLinearProgress-bar": {
    backgroundColor: "#3548D4",
  },
});

const StyledLink = styled(Link)({
  whiteSpace: "nowrap",
});

const StyledRow = styled(Grid)({
  marginBottom: "15px",
});

export interface BaseProjectCardProps {
  text: string;
  icon?: React.ReactNode;
  endAdornment?: React.ReactNode;
  loading?: boolean;
  error?: ClutchError | undefined;
}

interface ExtendedProjectCardProps extends BaseProjectCardProps {
  children?: React.ReactNode;
}

const LinkText = ({ text, link }: { text: string; link?: string }) => {
  const returnText = <Typography variant="body2">{text}</Typography>;

  if (link && text) {
    return <StyledLink href={link}>{returnText}</StyledLink>;
  }

  return returnText;
};

const TitleRow = ({ text, icon, endAdornment }: BaseProjectCardProps) => (
  <>
    {icon && (
      <Grid item xs={1}>
        {icon}
      </Grid>
    )}
    <Grid item xs={8}>
      <Typography variant="h4">{text}</Typography>
    </Grid>
    {endAdornment && (
      <Grid
        container
        item
        direction="row"
        xs={3}
        spacing={1}
        alignItems="center"
        justify="flex-end"
      >
        {endAdornment}
      </Grid>
    )}
  </>
);

const LastEvent = ({ time }: { time: number }) => (
  <>
    {time && (
      <>
        <Grid item>
          <FontAwesomeIcon icon={faClock} />
        </Grid>
        <Grid item>
          <Typography variant="body4">
            <EventTime date={setMilliseconds(time)} /> ago
          </Typography>
        </Grid>
      </>
    )}
  </>
);

const ProjectCard = ({
  children,
  text,
  icon,
  endAdornment,
  loading,
  error,
}: ExtendedProjectCardProps) => (
  <StyledCard container direction="row">
    <Grid container item direction="row" alignItems="flex-start">
      <TitleRow text={text} icon={icon} endAdornment={endAdornment} />
    </Grid>
    <StyledRow>
      {loading && (
        <StyledProgressContainer>
          {loading && <LinearProgress color="secondary" />}
        </StyledProgressContainer>
      )}
    </StyledRow>
    {children}
    {/* {error ? <Error subject={error} /> : children} */}
  </StyledCard>
);

export { LastEvent, LinkText, StyledCard, StyledRow, StyledLink };

export default ProjectCard;
