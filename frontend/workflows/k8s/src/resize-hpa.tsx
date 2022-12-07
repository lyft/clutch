import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import {
  Accordion,
  AccordionDetails,
  Button,
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
import _ from "lodash";
import { number, ref } from "yup";
import type Reference from "yup/lib/Reference";

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

  const oldHpaData = useDataLayout("oldHpaData");

  const metadataAnnotations = [];
  const metadataLabels = [];

  React.useEffect(() => {
    if (hpa.annotations) {
      _.forEach(hpa.annotations, (annotation, key) => {
        metadataAnnotations.push({ name: key, value: annotation });
      });
    }

    if (hpa.labels) {
      _.forEach(hpa.labels, (label, key) => {
        metadataLabels.push({ name: key, value: label });
      });
    }

    // save the original values of min and max replicas
    if (hpa) {
      oldHpaData.assign(hpa);
    }
  }, []);

  return (
    <WizardStep error={hpaData.error} isLoading={hpaData.isLoading}>
      <strong>HPA Details</strong>
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
            textFieldLabels: {
              disabledField: "Current min",
              updatedField: "New min",
            },
            input: {
              type: "number",
              key: "sizing.minReplicas",
              validation:
                hpa.sizing.minReplicas > 0
                  ? number().integer().moreThan(0)
                  : number().integer().min(0),
            },
          },
          {
            name: "Max Size",
            value: hpa.sizing.maxReplicas,
            textFieldLabels: {
              disabledField: "Current max",
              updatedField: "New max",
            },
            input: {
              type: "number",
              key: "sizing.maxReplicas",
              validation:
                hpa.sizing.minReplicas > 0
                  ? number()
                      .integer()
                      .min(ref("Min Size") as Reference<number>)
                  : number().integer().moreThan(0),
            },
          },
          { name: "Cluster", value: hpa.cluster },
        ]}
      />
      {metadataAnnotations.length > 0 && (
        <Accordion title="Annotations">
          <AccordionDetails>
            <MetadataTable data={metadataAnnotations} />
          </AccordionDetails>
        </Accordion>
      )}
      {metadataLabels.length > 0 && (
        <Accordion title="Labels">
          <AccordionDetails>
            <MetadataTable data={metadataLabels} />
          </AccordionDetails>
        </Accordion>
      )}
      <ButtonGroup>
        <Button text="Back" variant="neutral" onClick={() => onBack()} />
        <Button text="Resize" variant="destructive" onClick={onSubmit} />
      </ButtonGroup>
    </WizardStep>
  );
};

const Confirm: React.FC<ConfirmChild> = ({ notes }) => {
  const hpa = useDataLayout("hpaData").displayValue() as IClutch.k8s.v1.HPA;
  const resizeData = useDataLayout("resizeData");
  const oldHpaData = useDataLayout("oldHpaData").displayValue() as IClutch.k8s.v1.HPA;

  React.useEffect(() => {
  // if new values are either 50% bigger or smaller than old values, add a warning note
  const maxUpperBound = 1.5 * oldHpaData?.sizing?.maxReplicas;
  const maxLowerBound = 0.5 * oldHpaData?.sizing?.maxReplicas;
  const minUpperBound = 1.5 * oldHpaData?.sizing?.minReplicas;
  const minLowerBound = 0.5 * oldHpaData?.sizing?.minReplicas;
  const isMinReplicasDiffTooBig =
    hpa?.sizing?.minReplicas > minUpperBound || hpa?.sizing?.minReplicas < minLowerBound;
  const isMaxReplicasDiffTooBig =
    hpa?.sizing?.maxReplicas > maxUpperBound || hpa?.sizing?.maxReplicas < maxLowerBound;
  if (isMaxReplicasDiffTooBig || isMinReplicasDiffTooBig) {
    notes.unshift({
      text:
        "The new min or max size is more than 50% different from the old size. This may cause a large number of pods to be created or deleted.",
      severity: "warning",
    });
  }
  }, []);

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
          { name: "Old Min Size", value: oldHpaData?.sizing?.minReplicas },
          { name: "Old Max Size", value: oldHpaData?.sizing?.maxReplicas },
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
    oldHpaData: {},
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
