import React from "react";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import type { WizardChild } from "@clutch-sh/wizard";

import type { WorkflowProps } from ".";

const WelcomeStep: React.FC<WizardChild> = () => (
  <WizardStep isLoading={false} error="">Hello World!</WizardStep>
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
