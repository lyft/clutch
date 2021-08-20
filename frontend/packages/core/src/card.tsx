import * as React from "react";
import styled from "@emotion/styled";
import type {
  CardContentProps as MuiCardContentProps,
  CardHeaderProps as MuiCardHeaderProps,
} from "@material-ui/core";
import {
  Avatar,
  Card as MuiCard,
  CardActionArea,
  CardActionAreaProps,
  CardContent as MuiCardContent,
  Divider as MuiDivider,
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

const StyledCardHeader = styled.div({
  padding: "16px",
  display: "flex",
  alignItems: "center",
});

// TODO: add flex to this like the material component?
const StyledCardHeaderAvater = styled.div({
  marginRight: "16px",
});

const StyledCardHeaderSection = styled.div({
  paddingRight: "28px",
  paddingLeft: "28px",
});

// TODO: confirm the width
const Divider = styled(MuiDivider)({
  color: "rgba(13, 16, 48, 0.38)",
  height: "40px",
});

export interface CardHeaderSections {
  title: string;
  subheader?: string;
  titleColor?: string;
  subheaderColor?: string;
}

interface CardHeaderProps extends Pick<MuiCardHeaderProps, "avatar" | "title"> {
  children?: React.ReactNode;
  sections?: CardHeaderSections[];
}

// TODO: some flex/grid improvemnts
const CardHeader = ({ avatar, children, title, sections }: CardHeaderProps) => {
  return (
    <StyledCardHeaderContainer>
      <StyledCardHeader>
        <StyledCardHeaderAvater>{avatar}</StyledCardHeaderAvater>
        {/* todo wrap the title and sections together like material ui does? and add flex? */}
        <StyledTypography variant="h4" style={{ marginRight: "28px" }}>
          {title}
        </StyledTypography>
        {sections?.length > 0 &&
          sections.map(section => (
            <>
              <Divider orientation="vertical" flexItem />
              <StyledCardHeaderSection>
                <StyledTypography
                  variant="subtitle2"
                  color={section.titleColor ? section.titleColor : "#0D1030"}
                >
                  {section.title}
                </StyledTypography>
                {section.subheader && (
                  <StyledTypography
                    variant="body3"
                    color={
                      section.subheaderColor ? section.subheaderColor : "rgba(13, 16, 48, 0.6)"
                    }
                  >
                    {section.subheader}
                  </StyledTypography>
                )}
              </StyledCardHeaderSection>
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
