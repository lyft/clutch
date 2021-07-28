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
  CardHeader as MuiCardHeader,
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

interface CardHeaderProps extends Pick<MuiCardHeaderProps, "avatar" | "title"> {}

const CardHeader = ({ avatar, title }: CardHeaderProps) => (
  <MuiCardHeader
    style={{
      background: "#EBEDFB",
      padding: "16px",
    }}
    disableTypography
    avatar={avatar}
    title={<StyledTypography variant="h3">{title}</StyledTypography>}
  />
);

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
