import React from "react";
import { Error } from "@clutch-sh/core";
import { Container, Grid, Typography } from "@material-ui/core";
import styled from "styled-components";

const Heading = styled(Typography)`
  padding-left: 1.25rem;
`;

const Spacer = styled.div`
  margin: 30px;
`;

const SizedGrid = styled(Grid)`
  width: 100%;
  padding: 24px;
`;

interface PageLayoutProps {
  heading: string;
  error?: any;
}

const PageLayout: React.FC<PageLayoutProps> = ({ heading, error, children }) => {
  const hasError = error !== undefined && error !== "" && error !== null;
  return (
    <Spacer>
      {hasError && <Error message={error} />}
      <Container maxWidth="lg">
        <Heading variant="h5">
          <strong>{heading}</strong>
        </Heading>
        <SizedGrid>{children}</SizedGrid>
      </Container>
    </Spacer>
  );
};

export default PageLayout;
