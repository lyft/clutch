import React from "react";
import { Card as MuiCard, CardActionArea, CardContent } from "@material-ui/core";
import Typography from "@material-ui/core/Typography";
import styled from "styled-components";

const MIN_HEIGHT = "150px";

const SizedCard = styled(MuiCard)`
  ${({ theme, ...props }) => `
  min-height: ${MIN_HEIGHT};
  width: 300px;
  background: ${
    props["data-background"] !== undefined
      ? props["data-background"]
      : `linear-gradient(340deg, ${theme.palette.accent.main} 0%, ${theme.palette.secondary.main} 90%)`
  };
  `};
`;

const SizedContent = styled(CardContent)`
  min-height: ${MIN_HEIGHT};
`;

const CardTitle = styled(Typography)`
  ${({ ...props }) => `
    color: ${props["data-color"] !== undefined ? props["data-color"] : "#ffffff"};
  `};
`;

export interface CardProps {
  title: string;
  description: string;
  onClick?: () => void;
  titleColor?: string;
  descriptionColor?: "textPrimary" | "textSecondary";
  backgroundColor?: string;
}

const Card: React.FC<CardProps> = ({
  title,
  description,
  onClick,
  titleColor,
  descriptionColor = "textSecondary",
  backgroundColor,
}) => (
  <SizedCard raised data-background={backgroundColor}>
    <CardActionArea onClick={onClick}>
      <SizedContent>
        <CardTitle gutterBottom variant="h6" data-color={titleColor}>
          {title}
        </CardTitle>
        <Typography variant="body2" color={descriptionColor}>
          {description}
        </Typography>
      </SizedContent>
    </CardActionArea>
  </SizedCard>
);

export default Card;
