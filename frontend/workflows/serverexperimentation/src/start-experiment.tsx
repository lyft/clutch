import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import type { BaseWorkflowProps } from "@clutch-sh/core";
import {
  ButtonGroup,
  client,
  Confirmation,
  MetadataTable,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import * as yup from "yup";

const ClusterPairTargetDetails: React.FC<WizardChild> = () => {
  const { onSubmit } = useWizardContext();
  const clusterPairData = useDataLayout("clusterPairTargetData");
  const clusterPair = clusterPairData.displayValue();

  return (
    <WizardStep error={clusterPairData.error} isLoading={false}>
      <MetadataTable
        onUpdate={(key, value) => clusterPairData.updateData(key, value)}
        data={[
          {
            name: "Downstream Cluster",
            value: clusterPair.downstreamCluster,
            input: {
              key: "downstreamCluster",
              validation: yup.string().required(),
            },
          },
          {
            name: "Upstream Cluster",
            value: clusterPair.upstreamCluster,
            input: {
              key: "upstreamCluster",
              validation: yup.string().required(),
            },
          },
        ]}
      />
      <ButtonGroup
        buttons={[
          {
            text: "Next",
            onClick: onSubmit,
          },
        ]}
      />
    </WizardStep>
  );
};

const AbortExperimentDetails: React.FC<WizardChild> = () => {
  const { onSubmit, onBack } = useWizardContext();
  const abortExperimentData = useDataLayout("abortExperimentData");
  const abortExperiment = abortExperimentData.value;

  return (
    <WizardStep error={abortExperimentData.error} isLoading={false}>
      <MetadataTable
        onUpdate={(key, value) => abortExperimentData.updateData(key, value)}
        data={[
          {
            name: "Percent",
            value: abortExperiment.percent,
            input: {
              type: "number",
              key: "percent",
              validation: yup.number().integer().moreThan(-1).lessThan(101),
            },
          },
          {
            name: "HTTP Status",
            value: abortExperiment.httpStatus,
            input: {
              type: "number",
              key: "httpStatus",
              validation: yup.number().integer().moreThan(99).lessThan(600),
            },
          },
        ]}
      />
      <ButtonGroup
        buttons={[
          {
            text: "Back",
            onClick: onBack,
          },
          {
            text: "Next",
            onClick: onSubmit,
            destructive: true,
          },
        ]}
      />
    </WizardStep>
  );
};

const LatencyExperimentDetails: React.FC<WizardChild> = () => {
  const { onSubmit, onBack } = useWizardContext();
  const latencyExperimentData = useDataLayout("latencyExperimentData");
  const latencyExperiment = latencyExperimentData.value;

  return (
    <WizardStep error={latencyExperimentData.error} isLoading={false}>
      <MetadataTable
        onUpdate={(key, value) => latencyExperimentData.updateData(key, value)}
        data={[
          {
            name: "Percent",
            value: latencyExperiment.percent,
            input: {
              type: "number",
              key: "percent",
              validation: yup.number().integer().min(0).max(100),
            },
          },
          {
            name: "Duration (ms)",
            value: latencyExperiment.durationMs,
            input: {
              type: "number",
              key: "durationMs",
              validation: yup.number().integer().moreThan(0),
            },
          },
        ]}
      />
      <ButtonGroup
        buttons={[
          {
            text: "Back",
            onClick: onBack,
          },
          {
            text: "Next",
            onClick: onSubmit,
            destructive: true,
          },
        ]}
      />
    </WizardStep>
  );
};

const Confirm: React.FC<WizardChild> = () => {
  const startData = useDataLayout("startData");

  return (
    <WizardStep error={startData.error} isLoading={startData.isLoading}>
      <Confirmation action="Start" />
    </WizardStep>
  );
};

<<<<<<< HEAD:frontend/workflows/serverexperimentation/src/start-experiment.jsx
const createExperiment = data => {
  data["@type"] = "type.googleapis.com/clutch.chaos.serverexperimentation.v1.TestSpecification";
=======
const createExperiment = (data: IClutch.chaos.experimentation.v1.ITestSpecification) => {
>>>>>>> main:frontend/workflows/experimentation/src/start-experiment.tsx
  return client.post("/v1/experiments/create", {
    experiments: [
      {
        testConfig: data,
      },
    ],
  });
};

export const StartAbortExperiment: React.FC<BaseWorkflowProps> = ({ heading }) => {
  const dataLayout = {
    clusterPairTargetData: {},
    abortExperimentData: {},
    latencyExperimentData: {},
    startData: {
      deps: ["clusterPairTargetData", "abortExperimentData"],
      hydrator: (
        clusterPairTargetData: IClutch.chaos.experimentation.v1.IClusterPairTarget,
        abortExperimentData: IClutch.chaos.experimentation.v1.AbortFault
      ) => {
        return createExperiment({
          abort: {
            clusterPair: {
              downstreamCluster: clusterPairTargetData.downstreamCluster,
              upstreamCluster: clusterPairTargetData.upstreamCluster,
            },
            percent: abortExperimentData.percent,
            httpStatus: abortExperimentData.httpStatus,
          },
        });
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <ClusterPairTargetDetails name="target" />
      <AbortExperimentDetails name="abort" />
      <Confirm name="Confirmation" />
    </Wizard>
  );
};

export const StartLatencyExperiment: React.FC<BaseWorkflowProps> = ({ heading }) => {
  const dataLayout = {
    clusterPairTargetData: {},
    latencyExperimentData: {},
    startData: {
      deps: ["clusterPairTargetData", "latencyExperimentData"],
      hydrator: (
        clusterPairTargetData: IClutch.chaos.experimentation.v1.IClusterPairTarget,
        latencyExperimentData: IClutch.chaos.experimentation.v1.LatencyFault
      ) => {
        return createExperiment({
          latency: {
            clusterPair: {
              downstreamCluster: clusterPairTargetData.downstreamCluster,
              upstreamCluster: clusterPairTargetData.upstreamCluster,
            },
            percent: latencyExperimentData.percent,
            durationMs: latencyExperimentData.durationMs,
          },
        });
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <ClusterPairTargetDetails name="target" />
      <LatencyExperimentDetails name="latency" />
      <Confirm name="Confirmation" />
    </Wizard>
  );
};
