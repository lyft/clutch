import React from "react";
import type { ClutchError } from "@clutch-sh/core";
import { Error } from "@clutch-sh/core";
import styled from "@emotion/styled";
import { Container, Grid, Typography } from "@mui/material";

const PageContainer = styled.div({
  display: "flex",
  flex: "1 auto",
  margin: "30px",
});

const Heading = styled(Typography)({
  padding: "16px 0",
});

interface PageLayoutProps {
  heading: string;
  error?: ClutchError;
}

const PageLayout: React.FC<PageLayoutProps> = ({ heading, error, children }) => {
  const hasError = error !== undefined && error !== null;
  return (
    <PageContainer>
      <Container>
        <Heading variant="h5">
          <strong>{heading}</strong>
        </Heading>
        {hasError && <Error subject={error} />}
        <Grid>{children}</Grid>
      </Container>
    </PageContainer>
  );
};

export default PageLayout;
