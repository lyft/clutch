import React from "react";
import { Alert } from "@clutch-sh/core";
import styled from "@emotion/styled";
import { Container, Grid, Typography } from "@material-ui/core";

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
  error?: any;
}

const PageLayout: React.FC<PageLayoutProps> = ({ heading, error, children }) => {
  const hasError = error !== undefined && error !== "" && error !== null;
  return (
    <PageContainer>
      <Container>
        <Heading variant="h5">
          <strong>{heading}</strong>
        </Heading>
        {hasError && <Alert severity="error">{error}</Alert>}
        <Grid>{children}</Grid>
      </Container>
    </PageContainer>
  );
};

export default PageLayout;
