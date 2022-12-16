import React from "react";
import { CircularProgress, Grid, Paper } from "@mui/material";

import { styled } from "./Utils";

const LoadingSpinner = styled(CircularProgress)`
  color: #3548d4;
  position: absolute;
`;

const ContentContainer = styled(Grid)`
  position: relative;
`;

const ChildrenContainer = styled("div")({
  width: "100%",
});

const Overlay = styled(Paper)`
  position: absolute;
  height: 105%;
  width: 105%;
  display: flex;
  justify-content: center;
  align-items: center;
`;
const LoadingOveray = () => (
  <Overlay square elevation={0}>
    <LoadingSpinner />
  </Overlay>
);

interface LoadableProps {
  isLoading: boolean;
  variant?: "overlay";
}
const Loadable: React.FC<LoadableProps> = ({ isLoading, variant, children }) => {
  if (variant === "overlay") {
    return (
      <ContentContainer container direction="column" justifyContent="center" alignItems="center">
        <ChildrenContainer>{children}</ChildrenContainer>
        {isLoading && <LoadingOveray />}
      </ContentContainer>
    );
  }
  return isLoading ? <LoadingSpinner /> : <>{children}</>;
};

export default Loadable;
