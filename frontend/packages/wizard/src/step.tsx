import React from "react";
import { Error, Loadable, useWizardContext } from "@clutch-sh/core";
import { Grid } from "@material-ui/core";
import styled from "styled-components";

const SizedGrid = styled(Grid)`
  width: 100%;
  padding: 0 5%;
`;

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
    <SizedGrid container justify="center" direction="column" alignItems="center">
      {children}
    </SizedGrid>
  );
};

export default WizardStep;
