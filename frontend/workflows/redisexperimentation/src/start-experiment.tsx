import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
import type { BaseWorkflowProps, ClutchError } from "@clutch-sh/core";
import {
  Button,
  ButtonGroup,
  client,
  Dialog,
  DialogActions,
  DialogContent,
  Form,
} from "@clutch-sh/core";
import { PageLayout } from "@clutch-sh/experimentation";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";

import FormFields from "./form-fields";

enum FaultType {
  ERROR = "Error",
  LATENCY = "Latency",
}

interface RedisServiceCommandTargettingState {
  downstreamCluster: string;
  upstreamRedisCluster: string;
}

type ExperimentData = RedisServiceCommandTargettingState & {
  faultType: FaultType;
  percentage: number;
};

interface ExperimentDetailsProps {
  onStart: (ExperimentData) => void;
}

const ExperimentDetails: React.FC<ExperimentDetailsProps> = ({ onStart }) => {
  const initialExperimentData = {
    faultType: FaultType.ERROR,
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

  const fields = [
    {
      name: "downstreamCluster",
      label: "Downstream Cluster",
      type: "text",
      validation: yup.string().label("Downstream Cluster").required(),
      inputProps: { defaultValue: undefined },
    },
    {
      name: "upstreamRedisCluster",
      label: "Upstream Redis Cluster",
      type: "text",
      validation: yup.string().label("Upstream Redis Cluster").required(),
      inputProps: { defaultValue: undefined },
    },
    {
      name: "percentage",
      label: "Percent",
      type: "number",
      validation: yup.number().label("Percent").integer().min(1).max(100).required(),
      inputProps: { defaultValue: "0" },
    },
    {
      name: "faultType",
      label: "Fault Type",
      type: "select",
      inputProps: {
        options: [
          { label: "Error", value: FaultType.ERROR },
          { label: "Latency", value: FaultType.LATENCY },
        ],
        defaultValue: initialExperimentData.faultType,
      },
    },
  ];

  const schema: { [name: string]: yup.StringSchema | yup.NumberSchema } = {};
  fields
    .filter(field => field.validation !== undefined)
    .reduce((accumulator, field) => {
      if (field.validation) accumulator[field.name] = field.validation;
      return accumulator;
    }, schema);

  const { register, errors, handleSubmit } = useForm({
    mode: "onChange",
    reValidateMode: "onChange",
    resolver: yupResolver(yup.object().shape(schema)),
  });

  return (
    <Form onSubmit={handleSubmit(handleOnSubmit)}>
      <FormFields state={experimentDataState} items={fields} register={register} errors={errors} />
      <ButtonGroup>
        <Button text="Cancel" variant="neutral" onClick={handleOnCancel} />
        <Button text="Start" type="submit" />
      </ButtonGroup>
    </Form>
  );
};

const StartExperiment: React.FC<BaseWorkflowProps> = ({ heading }) => {
  const navigate = useNavigate();
  const [error, setError] = useState<ClutchError | undefined>(undefined);
  const [experimentData, setExperimentData] = useState<ExperimentData | undefined>(undefined);

  const handleOnCreatedExperiment = (id: number) => {
    navigate(`/experimentation/run/${id}`);
  };

  const handleOnCreatedExperimentFailure = (err: ClutchError) => {
    setExperimentData(undefined);
    setError(err);
  };

  const createExperiment = () => {
    const data = experimentData;
    const faultTargeting = {
      upstreamCluster: {
        name: data?.upstreamRedisCluster,
      },
      downstreamCluster: {
        name: data?.downstreamCluster,
      },
    } as IClutch.chaos.redisexperimentation.v1.FaultTargeting;

    const isError = data?.faultType === FaultType.ERROR;
    const fault = isError
      ? {
          errorFault: {
            percentage: {
              percentage: data?.percentage,
            },
          },
        }
      : {
          latencyFault: {
            percentage: {
              percentage: data?.percentage,
            },
          },
        };

    const testConfig = {
      "@type": "type.googleapis.com/clutch.chaos.redisexperimentation.v1.FaultConfig",
      faultTargeting,
      ...fault,
    };

    return client
      .post("/v1/chaos/experimentation/createExperiment", {
        data: {
          config: testConfig,
        },
      })
      .then(response => {
        handleOnCreatedExperiment(response?.data.experiment.runId);
      })
      .catch((err: ClutchError) => {
        handleOnCreatedExperimentFailure(err);
      });
  };

  return (
    <PageLayout heading={heading} error={error}>
      <ExperimentDetails onStart={experimentDetails => setExperimentData(experimentDetails)} />
      <Dialog
        title="Experiment Start Confirmation"
        open={experimentData !== undefined}
        onClose={() => setExperimentData(undefined)}
      >
        <DialogContent>
          Are you sure you want to start an experiment? The experiment will start immediately and
          you will be moved to experiment details view page.
        </DialogContent>
        <DialogActions>
          <Button text="No" variant="neutral" onClick={() => setExperimentData(undefined)} />
          <Button text="Yes" onClick={createExperiment} />
        </DialogActions>
      </Dialog>
    </PageLayout>
  );
};

export default StartExperiment;
