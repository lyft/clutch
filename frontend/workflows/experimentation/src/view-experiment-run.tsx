import React, { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { clutch, clutch as IClutch } from "@clutch-sh/api";
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
    IClutch.chaos.experimentation.v1.ExperimentRunDetails | undefined
  >(undefined);
  const [error, setError] = useState("");

  const { runID } = useParams();
  const navigate = useNavigate();

  function makeButtons() {
    const goBack = () => {
      navigate("/experimentation/list");
    };
    const goBackButton = {
      text: "Return",
      onClick: goBack,
    };

    const statusValue = clutch.chaos.experimentation.v1.Status[experiment.status].toString();
    if (statusValue === clutch.chaos.experimentation.v1.Status.COMPLETED.toString()) {
      return [goBackButton];
    }

    const title =
      statusValue === clutch.chaos.experimentation.v1.Status.RUNNING.toString()
        ? "Stop Experiment Run"
        : "Cancel Experiment Run";
    const destructiveButton = {
      text: title,
      destructive: true,
      onClick: () => {
        client
          .post("/v1/experiments/stop", { ids: [runID] })
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

  if (experiment === undefined && error === "") {
    client
      .post("/v1/experiments/details/run", { id: runID })
      .then(response => {
        setExperiment(response?.data?.runDetails);
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
          {experiment.properties.items.map(property => (
            <TextField
              key={property.label}
              label={property.label}
              defaultValue={property.value}
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
