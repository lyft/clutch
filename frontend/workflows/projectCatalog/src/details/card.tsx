import React from "react";
import { Card, Grid, Link, styled, Typography } from "@clutch-sh/core";
import { faClock } from "@fortawesome/free-regular-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

import { EventTime, setMilliseconds } from "./helpers";

const StyledCard = styled(Card)({
  width: "100%",
  height: "fit-content",
  padding: "15px",
});

const StyledLink = styled(Link)({
  whiteSpace: "nowrap",
});

const StyledRow = styled(Grid)({
  marginBottom: "15px",
});

export interface TitleRowProps {
  text: string;
  icon?: React.ReactNode;
  endAdornment?: React.ReactNode;
}

interface ProjectCardProps extends TitleRowProps {
  children?: React.ReactNode;
}

const LinkText = ({ text, link }: { text: string; link?: string }) => {
  const returnText = <Typography variant="body2">{text}</Typography>;

  if (link && text) {
    return <StyledLink href={link}>{returnText}</StyledLink>;
  }

  return returnText;
};

const TitleRow = ({ text, icon, endAdornment }: TitleRowProps) => (
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
    <Grid item>
      <FontAwesomeIcon icon={faClock} />
    </Grid>
    <Grid item>
      <Typography variant="body4">
        <EventTime date={setMilliseconds(time)} /> ago
      </Typography>
    </Grid>
  </>
);

const ProjectCard = ({ children, text, icon, endAdornment }: ProjectCardProps) => (
  <StyledCard container direction="row">
    <StyledRow container item direction="row" alignItems="flex-start">
      <TitleRow text={text} icon={icon} endAdornment={endAdornment} />
    </StyledRow>
    {children}
  </StyledCard>
);

export { LastEvent, LinkText, StyledCard, StyledRow, StyledLink };

export default ProjectCard;
