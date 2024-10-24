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
  NoteConfig,
  NotePanel,
  Resolver,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import _ from "lodash";
import { number, ref } from "yup";
import type Reference from "yup/lib/Reference";

import type { ConfirmChild, ResolverChild, WorkflowProps } from ".";

const WARNING_THRESHOLD = 0.5;

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

// The same notes will be displayed at both the final and next-to-last steps.
const HPADetails: React.FC<ConfirmChild> = ({ notes }) => {
  const { onSubmit, onBack } = useWizardContext();
  const hpaData = useDataLayout("hpaData");
  const hpa = hpaData.displayValue() as IClutch.k8s.v1.HPA;
  const update = (key: string, value: any) => {
    hpaData.updateData(key, value);
  };

  const currentHpaData = useDataLayout("currentHpaData");

  const metadataAnnotations = [];
  const metadataLabels = [];

  const [warnings, setWarnings] = React.useState<{
    minSize: NoteConfig;
    maxSize: NoteConfig;
  }>({ minSize: undefined, maxSize: undefined });
  const getWarning = React.useCallback(
    (current: number, newValue: number): boolean =>
      current * (1 / WARNING_THRESHOLD) < newValue || current / (1 / WARNING_THRESHOLD) > newValue,
    []
  );

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
      currentHpaData.assign(hpa);
    }
  }, []);

  React.useEffect(() => {
    if ((currentHpaData.displayValue() as IClutch.k8s.v1.HPA).sizing) {
      const {
        minReplicas,
        maxReplicas,
      } = (currentHpaData.displayValue() as IClutch.k8s.v1.HPA).sizing;
      setWarnings({
        minSize: getWarning(minReplicas, hpa.sizing.minReplicas)
          ? {
              severity: "warning",
              text: `New Min Size is more than ${WARNING_THRESHOLD * 100}% different`,
            }
          : undefined,
        maxSize: getWarning(maxReplicas, hpa.sizing.maxReplicas)
          ? {
              severity: "warning",
              text: `New Max Size is more than ${WARNING_THRESHOLD * 100}% different`,
            }
          : undefined,
      });
    }
  }, [hpa.sizing.minReplicas, hpa.sizing.maxReplicas]);

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
              warning: warnings.minSize,
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
              warning: warnings.maxSize,
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
        <Button text="Execute" variant="destructive" onClick={onSubmit} />
      </ButtonGroup>
      <NotePanel notes={notes} />
    </WizardStep>
  );
};

const Confirm: React.FC<ConfirmChild> = ({ notes }) => {
  const hpa = useDataLayout("hpaData").displayValue() as IClutch.k8s.v1.HPA;
  const resizeData = useDataLayout("resizeData");
  const currentHpaData = useDataLayout("currentHpaData").displayValue() as IClutch.k8s.v1.HPA;

  React.useEffect(() => {
    // if new values are either 50% bigger or smaller than old values, add a warning note
    const { maxReplicas, minReplicas } = currentHpaData.sizing;
    const maxUpperBound = (1 + WARNING_THRESHOLD) * maxReplicas;
    const maxLowerBound = (1 - WARNING_THRESHOLD) * maxReplicas;
    const minUpperBound = (1 + WARNING_THRESHOLD) * minReplicas;
    const minLowerBound = (1 - WARNING_THRESHOLD) * minReplicas;

    const isMinReplicasDiffTooBig =
      hpa.sizing.minReplicas > minUpperBound || hpa.sizing.minReplicas < minLowerBound;
    const isMaxReplicasDiffTooBig =
      hpa.sizing.maxReplicas > maxUpperBound || hpa.sizing.maxReplicas < maxLowerBound;
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
          { name: "Old Min Size", value: currentHpaData.sizing.minReplicas },
          { name: "Old Max Size", value: currentHpaData.sizing.maxReplicas },
          { name: "New Min Size", value: hpa.sizing.minReplicas },
          { name: "New Max Size", value: hpa.sizing.maxReplicas },
        ]}
      />
      <NotePanel notes={notes} />
    </WizardStep>
  );
};

const ResizeHPA: React.FC<WorkflowProps> = ({ resolverType, notes = [] }) => {
  const dataLayout = {
    hpaData: {},
    inputData: {},
    currentHpaData: {},
    resizeData: {
      deps: ["hpaData", "inputData", "currentHpaData"],
      hydrator: (
        hpaData: IClutch.k8s.v1.HPA,
        inputData: { clientset: string },
        currentHpaData: IClutch.k8s.v1.HPA
      ) => {
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
          currentSizing: {
            min: currentHpaData.sizing.minReplicas,
            max: currentHpaData.sizing.maxReplicas,
          },
        } as IClutch.k8s.v1.IResizeHPARequest);
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout}>
      <HPAIdentifier name="Lookup" resolverType={resolverType} />
      <HPADetails name="Modify" notes={notes} />
      <Confirm name="Result" notes={notes} />
    </Wizard>
  );
};

export default ResizeHPA;
