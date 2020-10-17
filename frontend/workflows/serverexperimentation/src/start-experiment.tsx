import React from "react";
import { clutch as IClutch } from "@clutch-sh/api";
import type { BaseWorkflowProps } from "@clutch-sh/core";
import {
  ButtonGroup,
  client,
  Confirmation,
  MetadataTable,
  Select,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import { FormControl, FormControlLabel, FormLabel, Radio, RadioGroup } from "@material-ui/core";
import * as yup from "yup";

interface RadioControlItem {
  label: string;
  value: string;
}

interface RadioControlProps {
  name: string;
  label: string;
  items: RadioControlItem[];
  onChange: (value: string) => void;
}

const RadioControl: React.FC<RadioControlProps> = ({ name, label, items, onChange }) => {
  return (
    <FormControl key={name}>
      <FormLabel component="legend">Upstream Cluster Type</FormLabel>
      <RadioGroup
        aria-label={label}
        name={name}
        defaultValue={items[0].value}
        onChange={e => onChange(e.target.value)}
      >
        {items &&
          items.map(item => {
            return (
              <FormControlLabel
                key={item.value}
                value={item.value}
                control={<Radio />}
                label={item.label}
              />
            );
          })}
      </RadioGroup>
    </FormControl>
  );
};

const faultInjectionTypeItems = [
  {
    label: "Internal",
    value: IClutch.chaos.serverexperimentation.v1.FaultInjectionType.FAULTINJECTIONTYPE_INGRESS.toString(),
  },
  {
    label: "External",
    value: IClutch.chaos.serverexperimentation.v1.FaultInjectionType.FAULTINJECTIONTYPE_EGRESS.toString(),
  },
];

const ClusterPairTargetDetails: React.FC<WizardChild> = () => {
  const { onSubmit } = useWizardContext();
  const clusterPairData = useDataLayout("clusterPairTargetData");
  const clusterPair = clusterPairData.displayValue();

  return (
    <WizardStep error={clusterPairData.error} isLoading={false}>
      <MetadataTable
        onUpdate={(key, value: string) => clusterPairData.updateData(key, value)}
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
      <RadioControl
        name="upstream_service_type"
        label="Upstream Service Type"
        items={faultInjectionTypeItems}
        onChange={(value: string) =>
          clusterPairData.updateData("faultInjectionType", parseInt(value, 10))
        }
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

enum FaultType {
  ABORT = "Abort",
  LATENCY = "Latency",
}

const ExperimentDetails: React.FC<WizardChild> = () => {
  const { onSubmit, onBack } = useWizardContext();
  const experimentData = useDataLayout("experimentData");
  const experiment = experimentData.value;

  const isAbort = (experiment?.type ?? FaultType.ABORT) === FaultType.ABORT;
  return (
    <WizardStep error={experimentData.error} isLoading={false}>
      <Select
        name="Type"
        label="Fault Type"
        options={[
          { label: "Abort", value: FaultType.ABORT },
          { label: "Latency", value: FaultType.LATENCY },
        ]}
        onChange={value => experimentData.updateData("type", value)}
      />
      <MetadataTable
        onUpdate={(key, value) => experimentData.updateData(key, value)}
        data={[
          {
            name: "Percent",
            value: experiment.percent,
            input: {
              type: "number",
              key: "percent",
              validation: yup.number().integer().min(1).max(100),
            },
          },
          isAbort
            ? {
                name: "HTTP Status",
                value: experiment.httpStatus,
                input: {
                  type: "number",
                  key: "httpStatus",
                  validation: yup.number().integer().min(100).max(599),
                },
              }
            : {
                name: "Duration (ms)",
                value: experiment.durationMs,
                input: {
                  type: "number",
                  key: "durationMs",
                  validation: yup.number().integer().min(1),
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

const StartExperiment: React.FC<BaseWorkflowProps> = ({ heading }) => {
  const createExperiment = (data: IClutch.chaos.serverexperimentation.v1.ITestConfig) => {
    const testConfig = data;
    testConfig["@type"] = "type.googleapis.com/clutch.chaos.serverexperimentation.v1.TestConfig";

    return client.post("/v1/chaos/experimentation/createExperiment", {
      config: testConfig,
    });
  };

  const dataLayout = {
    clusterPairTargetData: {},
    experimentData: {},
    startData: {
      deps: ["clusterPairTargetData", "experimentData"],
      hydrator: (
        clusterPairTargetData: IClutch.chaos.serverexperimentation.v1.IClusterPairTarget & {
          faultInjectionType: IClutch.chaos.serverexperimentation.v1.FaultInjectionType;
        },
        experimentData: IClutch.chaos.serverexperimentation.v1.AbortFaultConfig &
          IClutch.chaos.serverexperimentation.v1.LatencyFaultConfig & { type: FaultType }
      ) => {
        const isAbort = experimentData.type === FaultType.ABORT;
        const fault = isAbort
          ? { abort: { httpStatus: experimentData.httpStatus, percent: experimentData.percent } }
          : { latency: { durationMs: experimentData.durationMs, percent: experimentData.percent } };

        return createExperiment({
          clusterPair: {
            downstreamCluster: clusterPairTargetData.downstreamCluster,
            upstreamCluster: clusterPairTargetData.upstreamCluster,
          },
          faultInjectionType: clusterPairTargetData.faultInjectionType,
          ...fault,
        });
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <ClusterPairTargetDetails name="Target" />
      <ExperimentDetails name="Experiment Data" />
      <Confirm name="Confirmation" />
    </Wizard>
  );
};

export default StartExperiment;
