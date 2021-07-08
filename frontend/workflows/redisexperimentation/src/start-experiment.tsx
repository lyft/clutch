import React, { useState } from "react";
import { useForm } from "react-hook-form";
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
  useNavigate,
} from "@clutch-sh/core";
import { FormFields, PageLayout } from "@clutch-sh/experimentation";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";

enum FaultType {
  ERROR = "Error",
  LATENCY = "Latency",
}

type ExperimentData = {
  downstreamCluster: string;
  upstreamRedisCluster: string;
  environmentValue: string;
  faultType: FaultType;
  percentage: number;
};

interface ExperimentDetailsProps {
  environments: Environment[];
  onStart: (ExperimentData) => void;
}

const ExperimentDetails: React.FC<ExperimentDetailsProps> = ({ environments, onStart }) => {
  const initialExperimentData = {
    faultType: FaultType.ERROR,
    environmentValue: environments.length > 0 ? environments[0].value : "",
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
      name: "environmentValue",
      label: "Environment",
      type: "select",
      inputProps: {
        options: environments.map(env => {
          return {
            label: env.label,
            value: env.value,
          };
        }),
        defaultValue: initialExperimentData.environmentValue,
      },
      visible: environments.length > 0,
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

type Environment = {
  // If provided, it's displayed to the user instead of environment's value i.e. 'Staging' or 'Production'.
  label?: string;
  // The value that represents the environment i.e. 'staging' or 'production'.
  value: string;
};

interface StartExperimentProps extends BaseWorkflowProps {
  // The template that's used to resolve the final name of an upstream redis cluster i.e., "$[CLUSTER]-$[ENVIRONMENT]".
  // It defaults to '$[CLUSTER]' and supports the following variables:
  //  1. $CLUSTER - the value that a user entered in upstream cluster input field.
  //  2. $ENVIRONMENT - the value of environment that a user selected. It's available only if
  //                    'environments' is provided. Resolves to an empty string otherwise.
  upstreamRedisClusterTemplate?: string;
  // The template that's used to resolve the final name of a downstream cluster i.e., "$[CLUSTER]-$[ENVIRONMENT]".
  // It defaults to '$[CLUSTER]' and supports the following variables:
  //  1. $CLUSTER - the value that a user entered in upstream cluster input field.
  //  2. $ENVIRONMENT - the value of environment that a user selected. It's available only if
  //                    'environments' is provided. Resolves to an empty string otherwise.
  downstreamClusterTemplate?: string;
  // The list of environments that a user should be able to choose from
  environments?: Environment[];
}

const StartExperiment: React.FC<StartExperimentProps> = ({
  heading,
  upstreamRedisClusterTemplate,
  downstreamClusterTemplate,
  environments,
}) => {
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

  const evaluateClusterTemplate = (
    clusterTemplate: string,
    cluster: string,
    environment: string
  ) => {
    return clusterTemplate.replace("$[CLUSTER]", cluster).replace("$[ENVIRONMENT]", environment);
  };

  const createExperiment = (data: ExperimentData) => {
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
      <ExperimentDetails
        environments={environments ?? []}
        onStart={experimentDetails => setExperimentData(experimentDetails)}
      />
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
          <Button
            text="Yes"
            onClick={() => {
              const environment = experimentData.environmentValue;
              experimentData.upstreamRedisCluster = evaluateClusterTemplate(
                upstreamRedisClusterTemplate ?? "$[CLUSTER]",
                experimentData.upstreamRedisCluster,
                environment
              );
              experimentData.downstreamCluster = evaluateClusterTemplate(
                downstreamClusterTemplate ?? "$[CLUSTER]",
                experimentData.downstreamCluster,
                environment
              );
              createExperiment(experimentData);
            }}
          />
        </DialogActions>
      </Dialog>
    </PageLayout>
  );
};

export { StartExperiment, ExperimentDetails };
