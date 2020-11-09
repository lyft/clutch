import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { clutch as IClutch } from "@clutch-sh/api";
import type { BaseWorkflowProps } from "@clutch-sh/core";
import { ButtonGroup, client } from "@clutch-sh/core";
import { PageLayout } from "@clutch-sh/experimentation";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";

import Dialog from "./dialog";
import { FormContent } from "./form-content";

enum FaultType {
  ABORT = "Abort",
  LATENCY = "Latency",
}

const faultInjectionTypeItems = [
  {
    label: "Internal",
    value: IClutch.chaos.serverexperimentation.v1.FaultInjectionCluster.FAULTINJECTIONCLUSTER_UPSTREAM.toString(),
  },
  {
    label: "External (3rd party)",
    value: IClutch.chaos.serverexperimentation.v1.FaultInjectionCluster.FAULTINJECTIONCLUSTER_DOWNSTREAM.toString(),
  },
];

type ExperimentData = IClutch.chaos.serverexperimentation.v1.AbortFaultConfig &
  IClutch.chaos.serverexperimentation.v1.LatencyFaultConfig &
  IClutch.chaos.serverexperimentation.v1.ClusterPairTarget & { type: FaultType };

interface ExperimentDetailsProps {
  upstreamClusterTypeSelectionEnabled: boolean;
  onStart: (ExperimentData) => void;
}

const ExperimentDetails: React.FC<ExperimentDetailsProps> = ({
  upstreamClusterTypeSelectionEnabled,
  onStart,
}) => {
  const experimentDataState = useState<ExperimentData>({} as ExperimentData);
  const experimentData = experimentDataState[0];
  const navigate = useNavigate();

  const handleOnCancel = () => {
    navigate("/experimentation/list");
  };

  const handleOnSubmit = () => {
    onStart(experimentData);
  };

  const isAbort = (experimentData?.type ?? FaultType.ABORT) === FaultType.ABORT;
  const fields = [
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
    upstreamClusterTypeSelectionEnabled && {
      name: "upstreamClusterType",
      label: "Upstream Cluster Type",
      type: "radio-group",
      inputProps: { options: faultInjectionTypeItems },
    },
    {
      name: "type",
      label: "Fault Type",
      type: "select",
      inputProps: {
        options: [
          { label: "Abort", value: FaultType.ABORT },
          { label: "Latency", value: FaultType.LATENCY },
        ],
      },
    },
    {
      name: "percent",
      label: "Percent",
      type: "number",
      validation: yup.number().label("Percent").integer().min(1).max(100).required(),
      inputProps: { defaultValue: "0" },
    },
    isAbort
      ? {
          name: "httpStatus",
          label: "HTTP Status",
          type: "number",
          validation: yup.number().label("HTTP status").integer().min(100).max(599).required(),
          inputProps: { defaultValue: experimentData.httpStatus?.toString() },
        }
      : {
          name: "durationMs",
          label: "Duration (ms)",
          type: "number",
          validation: yup.number().label("Duration (ms)").integer().min(1).required(),
          inputProps: { defaultValue: experimentData.durationMs?.toString() },
        },
  ];

  const schema: { [name: string]: yup.StringSchema | yup.NumberSchema } = {};
  fields
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
      <FormContent state={experimentDataState} items={fields} register={register} errors={errors} />
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
}

const StartExperiment: React.FC<StartExperimentProps> = ({
  heading,
  upstreamClusterTypeSelectionEnabled = false,
}) => {
  const navigate = useNavigate();
  const [error, setError] = useState(undefined);
  const [experimentData, setExperimentData] = useState<ExperimentData | undefined>(undefined);

  const handleOnCreatedExperiment = (id: number) => {
    navigate(`/experimentation/run/${id}`);
  };

  const createExperiment = (
    data: IClutch.chaos.serverexperimentation.v1.AbortFaultConfig &
      IClutch.chaos.serverexperimentation.v1.LatencyFaultConfig &
      IClutch.chaos.serverexperimentation.v1.ClusterPairTarget & { type: FaultType }
  ) => {
    const isAbort = data.type === FaultType.ABORT;
    const fault = isAbort
      ? { abort: { httpStatus: data.httpStatus, percent: data.percent } }
      : { latency: { durationMs: data.durationMs, percent: data.percent } };

    const faultInjectionCluster =
      data.faultInjectionCluster ||
      IClutch.chaos.serverexperimentation.v1.FaultInjectionCluster.FAULTINJECTIONCLUSTER_UPSTREAM;

    return client
      .post("/v1/chaos/experimentation/createExperiment", {
        config: {
          "@type": "type.googleapis.com/clutch.chaos.serverexperimentation.v1.TestConfig",
          clusterPair: {
            downstreamCluster: data.downstreamCluster,
            upstreamCluster: data.upstreamCluster,
            faultInjectionCluster,
          },
          ...fault,
        },
      })
      .then(response => {
        handleOnCreatedExperiment(response?.data.experiment.id);
      })
      .catch(err => {
        setError(err.response.statusText);
      });
  };

  return (
    <PageLayout heading={heading} error={error}>
      <ExperimentDetails
        upstreamClusterTypeSelectionEnabled={upstreamClusterTypeSelectionEnabled}
        onStart={experimentDetails => setExperimentData(experimentDetails)}
      />
      <Dialog
        title="Experiment Start Confirmation"
        content="Are you sure you want to start an experiment? The experiment will start immediately and you will be moved to experiment details view page."
        open={experimentData !== undefined}
        onClose={() => setExperimentData(undefined)}
        buttons={[
          {
            label: "Yes",
            onAction: () => {
              createExperiment(experimentData);
            },
          },
          { label: "No", onAction: () => setExperimentData(undefined) },
        ]}
      />
    </PageLayout>
  );
};

export default StartExperiment;
