import React from "react";

{{- if .IsWizardTemplate}}
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
{{- end}}

import type { WorkflowProps } from ".";

{{- if .IsWizardTemplate}}

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

{{- else}}

const HelloWorld: React.FC<WorkflowProps> = ({ heading }) => {
  return <h1>Hello World - {heading}</h1>;
};

{{- end}}

export default HelloWorld;
