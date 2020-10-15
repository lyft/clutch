import React, { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { clutch as IClutch } from "@clutch-sh/api";
import { ButtonGroup, client, Error, TextField } from "@clutch-sh/core";
import styled from "styled-components";

import { propertyToString } from "./property-helpers";

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

    const statusValue = IClutch.chaos.experimentation.v1.Experiment.Status[
      experiment.status
    ].toString();
    const completedStatuses = [
      IClutch.chaos.experimentation.v1.Experiment.Status.STATUS_RUNNING.toString(),
      IClutch.chaos.experimentation.v1.Experiment.Status.STATUS_SCHEDULED.toString(),
    ];

    if (completedStatuses.indexOf(statusValue) < 0) {
      return [goBackButton];
    }

    const title =
      statusValue === IClutch.chaos.experimentation.v1.Experiment.Status.STATUS_RUNNING.toString()
        ? "Stop Experiment Run"
        : "Cancel Experiment Run";
    const destructiveButton = {
      text: title,
      destructive: true,
      onClick: () => {
        client
          .post("/v1/chaos/experimentation/cancelExperimentRun", { id: runID })
          .then(() => {
            setExperiment(undefined);
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
      .post("/v1/chaos/experimentation/getExperimentRunDetails", { id: runID })
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
              defaultValue={propertyToString(property)}
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
