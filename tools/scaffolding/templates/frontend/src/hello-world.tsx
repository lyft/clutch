import React from "react";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";

import type { WorkflowProps } from ".";

const WelcomeStep: React.FC<WizardChild> = () => (
  <WizardStep isLoading={false} error={undefined}>
    Hello World!
  </WizardStep>
);

const HelloWorld: React.FC<WorkflowProps> = ({ heading }) => {
  const dataLayout = {};
  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <WelcomeStep name="Welcome" />
    </Wizard>
  );
};

export default HelloWorld;
