import * as React from "react";
import styled from "@emotion/styled";
import type { CardContentProps as MuiCardContentProps } from "@material-ui/core";
import {
  Avatar,
  Card as MuiCard,
  CardActionArea,
  CardActionAreaProps,
  CardContent as MuiCardContent,
  Divider,
  Grid,
} from "@material-ui/core";

import { StyledTypography } from "./typography";

const StyledCard = styled(MuiCard)({
  boxShadow: "0px 4px 6px rgba(53, 72, 212, 0.2)",
  border: "1px solid rgba(13, 16, 48, 0.1)",

  ".MuiCardContent-root": {
    padding: "32px",
    color: "#0D1030",
    fontSize: "16px",
  },

  ".MuiCardActionArea-root:hover": {
    backgroundColor: "#F5F6FD",
  },

  ".MuiCardActionArea-root:active": {
    backgroundColor: "#D7DAF6",
  },
});

const StyledCardContent = styled(MuiCardContent)({
  "> .MuiPaper-root": {
    border: "0",
    borderRadius: "0",
  },
});

export interface CardProps {
  children?: React.ReactNode | React.ReactNode[];
}

const Card = ({ children, ...props }: CardProps) => <StyledCard {...props}>{children}</StyledCard>;

const StyledCardHeaderContainer = styled.div({
  background: "#EBEDFB",
});

const StyledCardHeader = styled(Grid)({
  padding: "16px",
  minHeight: "72px",
  margin: "0px",
  width: "100%",
});

const StyledCardHeaderAvater = styled.div({
  marginRight: "16px",
  fontSize: "24px",
});

const StyledDivider = styled(Divider)({
  color: "rgba(13, 16, 48, 0.38)",
  height: "36px",
  alignSelf: "center",
});

const StyledGridItem = styled(Grid)({
  textAlign: "center",
});

export interface CardHeaderSections {
  title: React.ReactNode;
  subheader?: React.ReactNode;
}

interface CardHeaderProps {
  avatar: React.ReactNode;
  children?: React.ReactNode;
  sections?: CardHeaderSections[];
  title: React.ReactNode;
}

// TODO: do we want to move the sections into the Dash card and keep this core card with just avatar and title
const CardHeader = ({ avatar, children, title, sections }: CardHeaderProps) => {
  return (
    <StyledCardHeaderContainer>
      <StyledCardHeader container wrap="nowrap" alignItems="center" spacing={2}>
        {/* TODO: use avatar component per design doc */}
        <StyledCardHeaderAvater>{avatar}</StyledCardHeaderAvater>
        <Grid item xs>
          <StyledTypography variant="h4">{title}</StyledTypography>
        </Grid>
        {sections?.length > 0 &&
          sections.map(section => (
            <>
              <StyledDivider orientation="vertical" />
              <StyledGridItem item xs>
                {section.title}
                {section.subheader}
              </StyledGridItem>
            </>
          ))}
      </StyledCardHeader>
      {children}
    </StyledCardHeaderContainer>
  );
};

interface CardContentProps extends MuiCardContentProps {}

const CardContent = ({ children, ...props }: CardContentProps) => (
  <StyledCardContent {...props}>{children}</StyledCardContent>
);

const StyledLandingCard = styled(Card)({
  border: "none",

  "& .header": {
    display: "inline-flex",
    marginBottom: "16px",
    fontWeight: "bold",
    fontSize: "12px",
    lineHeight: "36px",
    color: "rgba(13, 16, 48, 0.6)",
  },

  "& .header .icon .MuiAvatar-root": {
    height: "36px",
    width: "36px",
    marginRight: "8px",
    color: "rgba(13, 16, 48, 0.38)",
    backgroundColor: "rgba(13, 16, 48, 0.12)",
  },
});

export interface LandingCardProps extends Pick<CardActionAreaProps, "onClick"> {
  group: string;
  title: string;
  description: string;
}

export const LandingCard = ({ group, title, description, onClick, ...props }: LandingCardProps) => (
  <StyledLandingCard {...props}>
    <CardActionArea onClick={onClick}>
      <CardContent>
        <div className="header">
          <div className="icon">
            <Avatar>{group.charAt(0)}</Avatar>
          </div>
          <span>{group}</span>
        </div>
        <div>
          <StyledTypography variant="h3">{title}</StyledTypography>
          <StyledTypography style={{ color: "rgba(13, 16, 48, 0.6)" }} variant="body2">
            {description}
          </StyledTypography>
        </div>
      </CardContent>
    </CardActionArea>
  </StyledLandingCard>
);

export { Card, CardContent, CardHeader };
