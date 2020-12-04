import * as React from "react";
import { Card as MuiCard, CardActionArea, CardActionAreaProps, CardContent as MuiCardContent, CardHeader, Typography } from "@material-ui/core";
import styled from "@emotion/styled";
import { Avatar } from "@material-ui/core";

const StyledCard = styled(MuiCard)({
  boxShadow: "0px 4px 6px rgba(53, 72, 212, 0.2)",

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

export type CardProps = ({
  children?: React.ReactNode | React.ReactNode[];
});

export const Card = ({
  children, ...props
}: CardProps) => (
  <StyledCard {...props}>
    {children}
  </StyledCard>
);

const StyledLandingCard = styled(Card)({
  "& .header": {
    display: "inline-flex",
    marginBottom: "16px",
    textTransform: "uppercase",
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

  "& .title": {
    fontSize: "20px",
    fontWeight: 700,
  },

  "& .description": {
    marginTop: "5px",
    color: "rgba(13, 16, 48, 0.6)",
  }
});

export interface LandingCardProps extends Pick<CardActionAreaProps, "onClick"> {
  group: string;
  title: string;
  description: string;
};

export const LandingCard = ({ group, title, description, onClick, ...props }: LandingCardProps) => (
  <StyledLandingCard {...props}>
    <CardActionArea onClick={onClick}>
      <MuiCardContent>
        <div className="header">
          <div className="icon">
            <Avatar>{group.charAt(0)}</Avatar>
          </div>
          <span>{group}</span>
        </div>
        <div>
          <Typography className="title">{title}</Typography>
          <Typography className="description">{description}</Typography>
        </div>
      </MuiCardContent>
    </CardActionArea>
  </StyledLandingCard>
);

export const CardContent = MuiCardContent;

export default Card;
