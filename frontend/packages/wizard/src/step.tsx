import React from "react";
import { Error, Loadable, useWizardContext } from "@clutch-sh/core";
import styled from "@emotion/styled";
import { Grid as MuiGrid } from "@material-ui/core";

const Grid = styled(MuiGrid)({
  width: "100%",
  "> *": {
    padding: "8px",
  },
});

export interface WizardStepProps {
  isLoading: boolean;
  error: string;
}

const WizardStep: React.FC<WizardStepProps> = ({ isLoading, error, children }) => {
  const wizardContext = useWizardContext();
  const hasError = error !== undefined && error !== "" && error !== null;
  const showLoading = !hasError && isLoading;
  React.useEffect(() => {
    wizardContext.setIsLoading(showLoading);
  }, [showLoading]);
  if (showLoading) {
    return <Loadable isLoading={isLoading}>{children}</Loadable>;
  }
  return hasError ? (
    <Error message={error} />
  ) : (
    <Grid container justify="center" direction="column" alignItems="stretch">
      {children}
    </Grid>
  );
};

export default WizardStep;
