import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
import type { BaseWorkflowProps } from "@clutch-sh/core";
import { Button, ButtonGroup, client, Dialog } from "@clutch-sh/core";
import { PageLayout } from "@clutch-sh/experimentation";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";

import FormFields from "./form-fields";

enum FaultType {
  ABORT = "Abort",
  LATENCY = "Latency",
}

enum TargetType {
  REQUESTS = "requests",
  HOSTS = "hosts",
}

enum UpstreamClusterType {
  INTERNAL = "internal",
  EXTERNAL = "external",
}

type ExperimentData = {
  downstreamCluster: string;
  upstreamCluster: string;
  upstreamClusterType: UpstreamClusterType;
  targetType: TargetType;
  requestsPercentage: number;
  hostsPercentage: number;
  faultType: FaultType;
  httpStatus: number;
  durationMs: number;
};

interface ExperimentDetailsProps {
  upstreamClusterTypeSelectionEnabled: boolean;
  hostsPercentageBasedTargeting: boolean;
  onStart: (ExperimentData) => void;
}

const ExperimentDetails: React.FC<ExperimentDetailsProps> = ({
  upstreamClusterTypeSelectionEnabled,
  hostsPercentageBasedTargeting,
  onStart,
}) => {
  const initialExperimentData = {
    upstreamClusterType: UpstreamClusterType.INTERNAL,
    faultType: FaultType.ABORT,
    targetType: TargetType.REQUESTS,
  } as ExperimentData;

  const experimentDataState = useState<ExperimentData>(initialExperimentData);

  const experimentData = experimentDataState[0];
  const navigate = useNavigate();

  const handleOnCancel = () => {
    navigate("/experimentation/list");
  };

  const handleOnSubmit = () => {
    onStart(experimentData);
  };

  const faultInjectionClusterRadioGroup = {
    name: "faultInjectionCluster",
    label: "Upstream Cluster Type",
    type: "radio-group",
    visible:
      upstreamClusterTypeSelectionEnabled && experimentData.targetType === TargetType.REQUESTS,
    inputProps: {
      options: [
        {
          label: "Internal",
          value: UpstreamClusterType.INTERNAL,
        },
        {
          label: "External (3rd party)",
          value: UpstreamClusterType.EXTERNAL,
        },
      ],
      defaultValue: initialExperimentData.upstreamClusterType,
      disabled: experimentData.targetType !== TargetType.REQUESTS,
    },
  };
  const fakeFaultInjectionClusterRadioGroup = { ...faultInjectionClusterRadioGroup };
  fakeFaultInjectionClusterRadioGroup.name = "fakeFaultInjectionCluster";
  fakeFaultInjectionClusterRadioGroup.visible =
    upstreamClusterTypeSelectionEnabled && experimentData.targetType === TargetType.HOSTS;

  const isAbort = experimentData.faultType === FaultType.ABORT;
  const fields = [
    {
      label: "Cluster Pair",
      type: "title",
    },
    {
      name: "downstreamCluster",
      label: "Downstream Cluster",
      type: "text",
      validation: yup.string().label("Downstream Cluster").required(),
      inputProps: { defaultValue: undefined },
    },
    {
      name: "upstreamCluster",
      label: "Upstream Cluster",
      type: "text",
      validation: yup.string().label("Upstream Cluster").required(),
      inputProps: { defaultValue: undefined },
    },
    faultInjectionClusterRadioGroup,
    fakeFaultInjectionClusterRadioGroup,
    {
      label: "Targeting",
      type: "title",
    },
    {
      name: "targetType",
      label: "Target Type",
      type: "select",
      visible: hostsPercentageBasedTargeting,
      inputProps: {
        options: [
          {
            label: "Requests",
            value: TargetType.REQUESTS,
          },
          {
            label: "Hosts",
            value: TargetType.HOSTS,
          },
        ],
        defaultValue: initialExperimentData.targetType,
      },
    },
    {
      name: "requestsPercentage",
      label: "Percentage of Requests Served by All Hosts",
      type: "number",
      validation: yup.number().label("Percentage").integer().min(1).max(100).required(),
      visible: experimentData.targetType === TargetType.REQUESTS,
      inputProps: { defaultValue: "0" },
    },
    {
      name: "hostsPercentage",
      label: "Percentage of Hosts",
      type: "number",
      validation: yup.number().label("Percentage").integer().min(1).max(100).required(),
      visible: experimentData.targetType === TargetType.HOSTS,
      inputProps: { defaultValue: "0" },
    },
    {
      label: "Faults",
      type: "title",
    },
    {
      name: "faultType",
      label: "Fault Type",
      type: "select",
      inputProps: {
        options: [
          { label: "Abort", value: FaultType.ABORT },
          { label: "Latency", value: FaultType.LATENCY },
        ],
        defaultValue: initialExperimentData.faultType,
      },
    },
    {
      name: "httpStatus",
      label: "HTTP Status",
      type: "number",
      validation: yup.number().label("HTTP Status").integer().min(100).max(599).required(),
      visible: isAbort,
      inputProps: { defaultValue: experimentData.httpStatus?.toString() },
    },
    {
      name: "durationMs",
      label: "Duration (ms)",
      type: "number",
      validation: yup.number().label("Duration (ms)").integer().min(1).required(),
      visible: !isAbort,
      inputProps: { defaultValue: experimentData.durationMs?.toString() },
    },
  ];

  const schema: { [name: string]: yup.StringSchema | yup.NumberSchema } = {};
  const visibleFields = fields.filter(field => field.visible !== false);
  visibleFields
    .filter(field => field.validation !== undefined)
    .reduce((accumulator, field) => {
      accumulator[field.name] = field.validation;
      return accumulator;
    }, schema);

  const { register, errors, handleSubmit } = useForm({
    mode: "onChange",
    reValidateMode: "onChange",
    resolver: yupResolver(yup.object().shape(schema)),
  });

  return (
    <form onSubmit={handleSubmit(handleOnSubmit)}>
      <FormFields
        state={experimentDataState}
        items={visibleFields}
        register={register}
        errors={errors}
      />
      <ButtonGroup
        buttons={[
          {
            text: "Cancel",
            onClick: () => {
              handleOnCancel();
            },
          },
          {
            text: "Start",
            type: "submit",
          },
        ]}
      />
    </form>
  );
};

interface StartExperimentProps extends BaseWorkflowProps {
  upstreamClusterTypeSelectionEnabled?: boolean;
  hostsPercentageBasedTargeting?: boolean;
}

const StartExperiment: React.FC<StartExperimentProps> = ({
  heading,
  upstreamClusterTypeSelectionEnabled = false,
  hostsPercentageBasedTargeting = false,
}) => {
  const navigate = useNavigate();
  const [error, setError] = useState(undefined);
  const [experimentData, setExperimentData] = useState<ExperimentData | undefined>(undefined);

  const handleOnCreatedExperiment = (id: number) => {
    navigate(`/experimentation/run/${id}`);
  };

  const handleOnCreatedExperimentFailure = (err: string) => {
    setExperimentData(undefined);
    setError(err);
  };

  const createExperiment = (data: ExperimentData) => {
    const isUpstreamEnforcing = data.upstreamClusterType === UpstreamClusterType.INTERNAL;
    const isTargetingRequests = data.targetType === TargetType.REQUESTS;
    const isTargetingHosts = data.targetType === TargetType.HOSTS;

    const faultTargeting = {} as IClutch.chaos.serverexperimentation.v1.FaultTargeting;
    if (isUpstreamEnforcing) {
      if (isTargetingRequests) {
        faultTargeting.upstreamEnforcing = {
          downstreamCluster: {
            name: data.downstreamCluster,
          },
          upstreamCluster: {
            name: data.upstreamCluster,
          },
        };
      } else {
        faultTargeting.upstreamEnforcing = {
          downstreamCluster: {
            name: data.downstreamCluster,
          },
          upstreamPartialSingleCluster: {
            name: data.upstreamCluster,
            clusterPercentage: {
              percentage: isTargetingHosts ? data.hostsPercentage : 100,
            },
          },
        };
      }
    } else {
      faultTargeting.downstreamEnforcing = {
        downstreamCluster: {
          name: data.downstreamCluster,
        },
        upstreamCluster: {
          name: data.upstreamCluster,
        },
      };
    }

    const isAbort = data.faultType === FaultType.ABORT;
    let abortFault: IClutch.chaos.serverexperimentation.v1.AbortFault;
    let latencyFault: IClutch.chaos.serverexperimentation.v1.AbortFault;
    if (isAbort) {
      abortFault = {
        abortStatus: {
          httpStatusCode: data.httpStatus,
        },
        percentage: {
          percentage: isTargetingRequests ? data.requestsPercentage : 100,
        },
      } as IClutch.chaos.serverexperimentation.v1.AbortFault;
    } else {
      latencyFault = {
        latencyDuration: {
          fixedDurationMs: data.durationMs,
        },
        percentage: {
          percentage: isTargetingRequests ? data.requestsPercentage : 100,
        },
      } as IClutch.chaos.serverexperimentation.v1.LatencyFault;
    }

    return client
      .post("/v1/chaos/experimentation/createExperiment", {
        config: {
          "@type": "type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig",
          faultTargeting,
          abortFault,
          latencyFault,
        },
      })
      .then(response => {
        handleOnCreatedExperiment(response?.data.experiment.id);
      })
      .catch(err => {
        handleOnCreatedExperimentFailure(err.response.statusText);
      });
  };

  return (
    <PageLayout heading={heading} error={error}>
      <ExperimentDetails
        upstreamClusterTypeSelectionEnabled={upstreamClusterTypeSelectionEnabled}
        hostsPercentageBasedTargeting={hostsPercentageBasedTargeting}
        onStart={experimentDetails => setExperimentData(experimentDetails)}
      />
      <Dialog
        title="Experiment Start Confirmation"
        content="Are you sure you want to start an experiment? The experiment will start immediately and you will be moved to experiment details view page."
        open={experimentData !== undefined}
        onClose={() => setExperimentData(undefined)}
      >
        <Button
          text="Yes"
          onClick={() => {
            createExperiment(experimentData);
          }}
        />
        <Button text="No" onClick={() => setExperimentData(undefined)} />
      </Dialog>
    </PageLayout>
  );
};

export { StartExperiment, ExperimentDetails };
