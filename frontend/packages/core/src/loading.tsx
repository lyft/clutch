import React from "react";
import { CircularProgress, Grid, Paper } from "@material-ui/core";
import styled from "styled-components";

const LoadingSpinner = styled(CircularProgress)`
  ${({ theme }) => `
  color: ${theme.palette.accent.main};
  position: absolute;
  `}
`;

const ContentContainer = styled(Grid)`
  position: relative;
`;

const Overlay = styled(Paper)`
  position: absolute;
  height: 105%;
  width: 100%;
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
      <ContentContainer container direction="column" justify="center" alignItems="center">
        {children}
        {isLoading && <LoadingOveray />}
      </ContentContainer>
    );
  }
  return isLoading ? <LoadingSpinner /> : <>{children}</>;
};

export default Loadable;
