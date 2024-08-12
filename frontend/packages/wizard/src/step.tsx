import React from "react";
import type { ClutchError } from "@clutch-sh/core";
import { Error, Grid as ClutchGrid, Loadable, styled, useWizardContext } from "@clutch-sh/core";

const Grid = styled(ClutchGrid)({
  width: "100%",
  "> *": {
    margin: "8px 0",
  },
});

export interface WizardStepProps {
  isLoading: boolean;
  isComplete?: boolean;
  error?: ClutchError;
}

const WizardStep: React.FC<WizardStepProps> = ({ isComplete, isLoading, error, children }) => {
  const wizardContext = useWizardContext();
  const hasError = error !== undefined && error !== null;
  const showLoading = !hasError && isLoading;
  const completed = isComplete && !isLoading && !hasError;
  React.useEffect(() => {
    wizardContext.setIsLoading(showLoading);
  }, [showLoading]);
  React.useEffect(() => {
    wizardContext.setHasError(hasError);
  }, [error]);
  React.useEffect(() => {
    wizardContext.setIsComplete(completed);
  }, [completed]);
  if (showLoading) {
    return <Loadable isLoading={isLoading}>{children}</Loadable>;
  }
  return (
    <Grid container justifyContent="center" direction="column" alignItems="stretch">
      {hasError ? <Error subject={error} /> : children}
    </Grid>
  );
};

export default WizardStep;
