import * as React from "react";
import styled from "@emotion/styled";
import type {
  CardContentProps as MuiCardContentProps,
  CardHeaderProps as MuiCardHeaderProps,
} from "@material-ui/core";
import {
  Card as MuiCard,
  CardContent as MuiCardContent,
  CardHeader as MuiCardHeader,
} from "@material-ui/core";

import { StyledTypography } from "../typography";

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

export { Card, CardContent, CardHeader };
