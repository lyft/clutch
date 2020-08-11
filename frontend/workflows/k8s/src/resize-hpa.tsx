import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import {
  ButtonGroup,
  client,
  Confirmation,
  MetadataTable,
  NotePanel,
  Resolver,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import * as yup from "yup";

import type { ConfirmChild, ResolverChild, WorkflowProps } from ".";

const HPAIdentifier: React.FC<ResolverChild> = ({ resolverType }) => {
  const { onSubmit } = useWizardContext();
  const hpaData = useDataLayout("hpaData");
  const inputData = useDataLayout("inputData");

  const onResolve = ({ results, input }) => {
    // Decide how to process results.
    hpaData.assign(results[0]);
    inputData.assign(input);
    onSubmit();
  };

  return <Resolver type={resolverType} searchLimit={1} onResolve={onResolve} />;
};

const HPADetails: React.FC<WizardChild> = () => {
  const { onSubmit, onBack } = useWizardContext();
  const hpaData = useDataLayout("hpaData");
  const hpa = hpaData.displayValue() as IClutch.k8s.v1.HPA;
  const update = (key: string, value: any) => {
    hpaData.updateData(key, value);
  };

  return (
    <WizardStep error={hpaData.error} isLoading={hpaData.isLoading}>
      <MetadataTable
        onUpdate={update}
        data={[
          { name: "Name", value: hpa.name },
          { name: "Namespace", value: hpa.namespace },
          { name: "Current Replicas", value: hpa.sizing.currentReplicas },
          { name: "Desired Replicas", value: hpa.sizing.desiredReplicas },
          {
            name: "Min Size",
            value: hpa.sizing.minReplicas,
            input: {
              type: "number",
              key: "sizing.minReplicas",
              validation: yup.number().integer().moreThan(0),
            },
          },
          {
            name: "Max Size",
            value: hpa.sizing.maxReplicas,
            input: { type: "number", key: "sizing.maxReplicas" },
          },
          { name: "Cluster", value: hpa.cluster },
        ]}
      />
      <ButtonGroup
        buttons={[
          {
            text: "Back",
            onClick: onBack,
          },
          {
            text: "Resize",
            onClick: onSubmit,
            destructive: true,
          },
        ]}
      />
    </WizardStep>
  );
};

const Confirm: React.FC<ConfirmChild> = ({ notes }) => {
  const hpa = useDataLayout("hpaData").displayValue() as IClutch.k8s.v1.HPA;
  const resizeData = useDataLayout("resizeData");

  return (
    <WizardStep error={resizeData.error} isLoading={resizeData.isLoading}>
      <Confirmation action="Resize" />
      <MetadataTable
        data={[
          { name: "Name", value: hpa.name },
          { name: "Namespace", value: hpa.namespace },
          { name: "Cluster", value: hpa.cluster },
          { name: "New Min Size", value: hpa.sizing.minReplicas },
          { name: "New Max Size", value: hpa.sizing.maxReplicas },
        ]}
      />
      <NotePanel notes={notes} />
    </WizardStep>
  );
};

const ResizeHPA: React.FC<WorkflowProps> = ({ heading, resolverType, notes = [] }) => {
  const dataLayout = {
    hpaData: {},
    inputData: {},
    resizeData: {
      deps: ["hpaData", "inputData"],
      hydrator: (hpaData: IClutch.k8s.v1.HPA, inputData: { clientset: string }) => {
        const clientset = inputData.clientset ?? "unspecified";

        return client.post("/v1/k8s/resizeHPA", {
          clientset,
          cluster: hpaData.cluster,
          namespace: hpaData.namespace,
          name: hpaData.name,
          sizing: {
            min: hpaData.sizing.minReplicas,
            max: hpaData.sizing.maxReplicas,
          },
        } as IClutch.k8s.v1.IResizeHPARequest);
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <HPAIdentifier name="Lookup" resolverType={resolverType} />
      <HPADetails name="Modify" />
      <Confirm name="Confirmation" notes={notes} />
    </Wizard>
  );
};

export default ResizeHPA;
