import React from "react";
import { Grid, GridJustification } from "@material-ui/core";
import FiberManualRecordTwoToneIcon from "@material-ui/icons/FiberManualRecordTwoTone";
import MuiTrendingUpIcon from "@material-ui/icons/TrendingUp";
import styled from "styled-components";

const StatusIcon = styled(FiberManualRecordTwoToneIcon)`
  ${({ ...props }) => `
    color: ${props["data-color"]}
  `}
`;

export interface StatusProps {
  variant?: "neutral" | "success" | "failure";
  align?: "left" | "center" | "right";
}

const Status: React.FC<StatusProps> = ({
  children,
  variant = "neutral",
  align = "left",
  ...props
}) => {
  let justifyContent: GridJustification = "flex-start";
  if (align === "right") {
    justifyContent = "flex-end";
  } else if (align === "center") {
    justifyContent = "center";
  }
  return (
    <Grid container alignItems="center" justify={justifyContent} {...props}>
      {variant === "neutral" && (
        <>
          <StatusIcon data-color="darkgray" /> {children}
        </>
      )}
      {variant === "success" && (
        <>
          <StatusIcon data-color="limegreen" /> {children}
        </>
      )}
      {variant === "failure" && (
        <>
          <StatusIcon data-color="red" /> {children}
        </>
      )}
    </Grid>
  );
};

const StyledTrendingUpIcon = styled(MuiTrendingUpIcon)`
  ${({ theme }) => `
  color: ${theme.palette.accent.main};
  margin: 1%;
  `}
`;

const TrendingUpIcon: React.FC = () => <StyledTrendingUpIcon />;

export { Status, TrendingUpIcon };
