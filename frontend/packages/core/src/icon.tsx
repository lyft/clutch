import React from "react";
import styled from "@emotion/styled";
import FiberManualRecordTwoToneIcon from "@mui/icons-material/FiberManualRecordTwoTone";
import { Grid } from "@mui/material";

import type { GridJustification } from "./grid";

const StyledStatusIcon = styled(FiberManualRecordTwoToneIcon)`
  ${({ ...props }) => `
    color: ${props["data-color"]}
  `}
`;

export const AvatarIcon = styled.img({
  borderRadius: "50%",
  maxHeight: "100%",
  maxWidth: "initial",
});

export interface StatusProps {
  variant?: "neutral" | "success" | "failure";
  align?: "left" | "center" | "right";
}

export const StatusIcon: React.FC<StatusProps> = ({
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
    <Grid container alignItems="center" justifyContent={justifyContent} {...props}>
      {variant === "neutral" && (
        <>
          <StyledStatusIcon data-color="#C2C8F2" /> {children}
        </>
      )}
      {variant === "success" && (
        <>
          <StyledStatusIcon data-color="#69F0AE" /> {children}
        </>
      )}
      {variant === "failure" && (
        <>
          <StyledStatusIcon data-color="#FF8A80" /> {children}
        </>
      )}
    </Grid>
  );
};
