import React from "react";
import type { ClutchError } from "@clutch-sh/core";
import { Error } from "@clutch-sh/core";
import { Grid } from "@mui/material";

interface PageLayoutProps {
  error?: ClutchError;
}

const PageLayout: React.FC<PageLayoutProps> = ({ error, children }) => {
  const hasError = error !== undefined && error !== null;
  return (
    <Grid spacing={2}>
      {hasError && <Error subject={error} />}
      <Grid>{children}</Grid>
    </Grid>
  );
};

export default PageLayout;
