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
`;
const LoadingOveray = () => (
  <Overlay square elevation={0}>
    <LoadingSpinner />
  </Overlay>
);

export interface LoadableProps {
  isLoading: boolean;
  overlay?: boolean;
}

const Loadable: React.FC<LoadableProps> = ({ isLoading, overlay = false, children }) => {
  if (overlay) {
    return (
      <ContentContainer container>
        {children}
        {isLoading && <LoadingOveray />}
      </ContentContainer>
    );
  }
  return isLoading ? <LoadingSpinner /> : <>{children}</>;
};

export default Loadable;
