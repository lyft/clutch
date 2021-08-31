import * as React from "react";
import styled from "@emotion/styled";
import type { CardHeaderProps as MuiCardHeaderProps } from "@material-ui/core";
import {
  Avatar,
  Card as MuiCard,
  CardActionArea,
  CardActionAreaProps,
  CardHeader as MuiCardHeader,
} from "@material-ui/core";
import { spacing } from "@material-ui/system";

import { StyledTypography } from "./typography";

const StyledCard = styled(MuiCard)({
  boxShadow: "0px 4px 6px rgba(53, 72, 212, 0.2)",
  border: "1px solid rgba(13, 16, 48, 0.1)",

  ".MuiCardContent-root": {
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

export interface CardProps {
  children?: React.ReactNode | React.ReactNode[];
}

const Card = ({ children, ...props }: CardProps) => <StyledCard {...props}>{children}</StyledCard>;

const StyledCardHeaderContainer = styled.div({
  background: "#EBEDFB",
});

interface CardHeaderProps extends Pick<MuiCardHeaderProps, "avatar" | "title"> {
  children?: React.ReactNode;
}

const CardHeader = ({ avatar, children, title }: CardHeaderProps) => (
  <StyledCardHeaderContainer>
    <MuiCardHeader
      style={{
        padding: "16px",
      }}
      disableTypography
      avatar={avatar}
      title={<StyledTypography variant="h3">{title}</StyledTypography>}
    />
    {children}
  </StyledCardHeaderContainer>
);

// Material UI Spacing system supports many props https://material-ui.com/system/spacing/#api
// We can add more to this list as use cases arise
interface SpacingProps {
  padding?: number;
  // shorthand for padding
  p?: number;
}

const BaseCardContent = styled.div<SpacingProps>`
  ${spacing}
`;

const StyledCardContent = styled(BaseCardContent)({
  "> .MuiPaper-root": {
    border: "0",
    borderRadius: "0",
  },
});

interface CardContentProps extends SpacingProps {
  children?: React.ReactNode | React.ReactNode[];
}

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
      <CardContent padding={4}>
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
