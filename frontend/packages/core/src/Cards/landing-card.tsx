import * as React from "react";
import styled from "@emotion/styled";
import { Avatar, CardActionArea, CardActionAreaProps } from "@material-ui/core";

import { StyledTypography } from "../typography";

import { Card, CardContent } from "./card";

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
