import React from "react";
import type { ClutchError } from "@clutch-sh/core";
import { Error, Loadable, styled, useWizardContext } from "@clutch-sh/core";
import { Grid as MuiGrid } from "@material-ui/core";

const Grid = styled(MuiGrid)({
  width: "100%",
  "> *": {
    margin: "8px 0",
  },
});

export interface WizardStepProps {
  isLoading: boolean;
  error: ClutchError;
}

const WizardStep: React.FC<WizardStepProps> = ({ isLoading, error, children }) => {
  const wizardContext = useWizardContext();
  const hasError = error !== undefined && error !== null;
  const showLoading = !hasError && isLoading;
  React.useEffect(() => {
    wizardContext.setIsLoading(showLoading);
  }, [showLoading]);
  React.useEffect(() => {
    wizardContext.setHasError(hasError);
  }, [error]);
  if (showLoading) {
    return <Loadable isLoading={isLoading}>{children}</Loadable>;
  }
  return (
    <Grid container justify="center" direction="column" alignItems="stretch">
      {hasError ? <Error subject={error} /> : children}
    </Grid>
  );
};

export default WizardStep;
