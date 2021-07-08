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

import type { FormItem } from "../../experimentation/src/core/form-fields";

enum FaultType {
  ABORT = "Abort",
  LATENCY = "Latency",
}

enum UpstreamClusterType {
  INTERNAL = "internal",
  EXTERNAL = "external",
}

type ExperimentData = {
  downstreamCluster: string;
  upstreamCluster: string;
  environmentValue: string;
  upstreamClusterType: UpstreamClusterType;
  requestsPercentage: number;
  hostsPercentage: number;
  faultType: FaultType;
  httpStatus: number;
  durationMs: number;
};

interface ExperimentDetailsProps {
  upstreamClusterTypeSelectionEnabled: boolean;
  environments: Environment[];
  onStart: (ExperimentData) => void;
}

const ExperimentDetails: React.FC<ExperimentDetailsProps> = ({
  upstreamClusterTypeSelectionEnabled,
  environments,
  onStart,
}) => {
  const initialExperimentData = {
    upstreamClusterType: UpstreamClusterType.INTERNAL,
    faultType: FaultType.ABORT,
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
    {
      name: "upstreamClusterType",
      label: "Upstream Cluster Type",
      type: "radio-group",
      visible: upstreamClusterTypeSelectionEnabled,
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
      },
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
      name: "requestsPercentage",
      label: "Percentage of Requests Served by All Hosts",
      type: "number",
      validation: yup.number().label("Percentage").integer().min(1).max(100).required(),
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
  ] as FormItem[];

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
    <Form onSubmit={handleSubmit(handleOnSubmit)}>
      <FormFields
        state={experimentDataState}
        items={visibleFields}
        register={register}
        errors={errors}
      />
      <ButtonGroup>
        <Button text="Cancel" variant="neutral" onClick={handleOnCancel} />
        <Button text="Start" type="submit" />
      </ButtonGroup>
    </Form>
  );
};

type Environment = {
  // If provided, it's displayed to the user instead of environment's value i.e., 'Staging' or 'Production'.
  label?: string;
  // The value that represents the environment i.e., 'staging' or 'production'.
  value: string;
};

// The set of templates for upstream cluster.
type UpstreamClusterTemplates = {
  // The template that's used to resolve the final name of an internal upstream cluster.
  // It supports the following variables:
  //  1. $CLUSTER - the value that a user entered in upstream cluster input field.
  //  2. $ENVIRONMENT - the value of environment that a user selected. It's available only if
  //                    'environments' is provided. Resolves to an empty string otherwise.
  internalClusterTemplate: string;
  // The template that's used to resolve the final name of an external upstream cluster.
  // It supports the following variables:
  //  1. $CLUSTER - the value that a user entered in upstream cluster input field.
  //  2. $ENVIRONMENT - the value of environment that a user selected. It's available only if
  //                    'environments' is provided. Resolves to an empty string otherwise.
  externalClusterTemplate?: string;
};

interface StartExperimentProps extends BaseWorkflowProps {
  // Templates that are used to resolve the final name of an upstream cluster i.e., "$[CLUSTER]-$[ENVIRONMENT]".
  // It supports providing separate templates for internal and external upstream clusters. Both default to
  // '$[CLUSTER]' and support the following variables:
  //  1. $CLUSTER - the value that a user entered in upstream cluster input field.
  //  2. $ENVIRONMENT - the value of environment that a user selected. It's available only if
  //                    'environments' is provided. Resolves to an empty string otherwise.
  upstreamClusterTemplates?: UpstreamClusterTemplates;
  // The templates that's used to resolve the final name of a downstream cluster i.e., "$[CLUSTER]-$[ENVIRONMENT]".
  // It defaults to '$[CLUSTER]' and supports the following variables:
  //  1. $CLUSTER - the value that a user entered in upstream cluster input field.
  //  2. $ENVIRONMENT - the value of environment that a user selected. It's available only if
  //                    'environments' is provided. Resolves to an empty string otherwise.
  downstreamClusterTemplate?: string;
  // Whether a user should be able to select if an upstream cluster is an external or internal dependency.
  upstreamClusterTypeSelectionEnabled?: boolean;
  // The list of environments that a user should be able to choose from. Defaults to empty list. If not present or empty,
  // a user is not asked for environment.
  environments?: Environment[];
}

const StartExperiment: React.FC<StartExperimentProps> = ({
  heading,
  upstreamClusterTemplates,
  downstreamClusterTemplate,
  upstreamClusterTypeSelectionEnabled,
  environments,
}) => {
  const navigate = useNavigate();
  const [error, setError] = useState(undefined);
  const [experimentData, setExperimentData] = useState<ExperimentData | undefined>(undefined);

  const handleOnCreatedExperiment = (id: string) => {
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
    const isUpstreamEnforcing = data.upstreamClusterType === UpstreamClusterType.INTERNAL;

    const faultTargeting = {} as IClutch.chaos.serverexperimentation.v1.FaultTargeting;
    if (isUpstreamEnforcing) {
      faultTargeting.upstreamEnforcing = {
        downstreamCluster: {
          name: data.downstreamCluster,
        },
        upstreamCluster: {
          name: data.upstreamCluster,
        },
      };
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
          percentage: data.requestsPercentage,
        },
      } as IClutch.chaos.serverexperimentation.v1.AbortFault;
    } else {
      latencyFault = {
        latencyDuration: {
          fixedDurationMs: data.durationMs,
        },
        percentage: {
          percentage: data.requestsPercentage,
        },
      } as IClutch.chaos.serverexperimentation.v1.LatencyFault;
    }

    return client
      .post("/v1/chaos/experimentation/createExperiment", {
        data: {
          config: {
            "@type": "type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig",
            faultTargeting,
            abortFault,
            latencyFault,
          },
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
        upstreamClusterTypeSelectionEnabled={upstreamClusterTypeSelectionEnabled ?? false}
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
          <Button variant="neutral" text="No" onClick={() => setExperimentData(undefined)} />
          <Button
            text="Yes"
            onClick={() => {
              const environment = experimentData.environmentValue;
              let upstreamClusterTemplate = upstreamClusterTemplates?.internalClusterTemplate;
              if (experimentData.upstreamClusterType === UpstreamClusterType.EXTERNAL) {
                upstreamClusterTemplate = upstreamClusterTemplates?.externalClusterTemplate;
              }
              experimentData.upstreamCluster = evaluateClusterTemplate(
                upstreamClusterTemplate ?? "$[CLUSTER]",
                experimentData.upstreamCluster,
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
