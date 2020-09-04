import React, { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { clutch, clutch as IClutch } from "@clutch-sh/api";
import type { ButtonProps } from "@clutch-sh/core";
import { ButtonGroup, client, Error, TextField } from "@clutch-sh/core";
import styled from "styled-components";

export const Form = styled.form`
  align-items: center;
  display: flex;
  flex-direction: column;
  width: 100%;
`;

const ViewExperimentRun: React.FC = () => {
  const [experiment, setExperiment] = useState<
    IClutch.chaos.experimentation.v1.ExperimentRunConfigPairDetails | undefined
  >(undefined);
  const [error, setError] = useState("");

  const { id } = useParams();
  const navigate = useNavigate();

  function makeButtons(): ButtonProps[] {
    const goBack = () => {
      navigate("/experimentation/list");
    };
    const goBackButton = {
      text: "Return",
      onClick: goBack,
    };

    if (experiment.status === clutch.chaos.experimentation.v1.Status.COMPLETED) {
      return [goBackButton];
    }

    const title =
      experiment.status === clutch.chaos.experimentation.v1.Status.RUNNING
        ? "Stop Experiment Run"
        : "Cancel Experiment Run";
    const destructiveButton = {
      text: title,
      destructive: true,
      onClick: () => {
        client
          .post("/v1/experiments/stop", { ids: [experiment.runId] })
          .then(() => {
            goBack();
          })
          .catch(err => {
            setError(err.response.statusText);
          });
      },
    };

    return [goBackButton, destructiveButton];
  }

  if (experiment === undefined) {
    client
      .post("/v1/experiment/details/run-config", { id })
      .then(response => {
        setExperiment(response?.data?.runConfigPairDetails);
      })
      .catch(err => {
        setError(err.response.statusText);
      });
  }

  return (
    <Form>
      {error && <Error message={error} />}
      {experiment && (
        <>
          {experiment.form.fields.map(field => (
            <TextField
              key={field.label}
              label={field.label}
              defaultValue={field.value}
              InputProps={{ readOnly: true }}
            />
          ))}
          <TextField
            multiline
            label="Config"
            defaultValue={JSON.stringify(experiment.config, null, 4)}
            InputProps={{ readOnly: true }}
          />
          <ButtonGroup buttons={makeButtons()} />
        </>
      )}
    </Form>
  );
};

export default ViewExperimentRun;
