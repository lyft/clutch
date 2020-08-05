import React from "react";
import {
  ButtonGroup,
  client,
  Confirmation,
  MetadataTable,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import * as yup from "yup";

const ClusterPairTargetDetails = () => {
  const { onSubmit } = useWizardContext();
  const clusterPairData = useDataLayout("clusterPairTargetData");
  const clusterPair = clusterPairData.displayValue();

  return (
    <WizardStep error={clusterPairData.error}>
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

const AbortExperimentDetails = () => {
  const { onSubmit, onBack } = useWizardContext();
  const abortExperimentData = useDataLayout("abortExperimentData");
  const abortExperiment = abortExperimentData.value;

  return (
    <WizardStep error={abortExperimentData.error}>
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

const LatencyExperimentDetails = () => {
  const { onSubmit, onBack } = useWizardContext();
  const latencyExperimentData = useDataLayout("latencyExperimentData");
  const latencyExperiment = latencyExperimentData.value;

  return (
    <WizardStep error={latencyExperimentData.error}>
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

const Confirm = () => {
  const startData = useDataLayout("startData");

  return (
    <WizardStep error={startData.error} isLoading={startData.isLoading}>
      <Confirmation action="Start" />
    </WizardStep>
  );
};

const createExperiment = data => {
  data["@type"] = "type.googleapis.com/clutch.chaos.serverexperimentation.v1.TestSpecification";
  return client.post("/v1/experiments/create", {
    experiments: [
      {
        testConfig: data,
      },
    ],
  });
};

export const StartAbortExperiment = ({ heading }) => {
  const dataLayout = {
    clusterPairTargetData: {},
    abortExperimentData: {},
    latencyExperimentData: {},
    startData: {
      deps: ["clusterPairTargetData", "abortExperimentData"],
      hydrator: (clusterPairTargetData, abortExperimentData) => {
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

export const StartLatencyExperiment = ({ heading }) => {
  const dataLayout = {
    clusterPairTargetData: {},
    latencyExperimentData: {},
    startData: {
      deps: ["clusterPairTargetData", "latencyExperimentData"],
      hydrator: (clusterPairTargetData, latencyExperimentData) => {
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
