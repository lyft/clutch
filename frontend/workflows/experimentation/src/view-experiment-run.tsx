import React, { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { clutch as IClutch } from "@clutch-sh/api";
import type { ClutchError } from "@clutch-sh/core";
import {
  BaseWorkflowProps,
  Button,
  ButtonGroup,
  client,
  Form,
  Link,
  TextField,
} from "@clutch-sh/core";

import PageLayout from "./core/page-layout";
import { propertyToString } from "./property-helpers";

const ViewExperimentRun: React.FC<BaseWorkflowProps> = ({ heading }) => {
  const [experiment, setExperiment] = useState<
    IClutch.chaos.experimentation.v1.ExperimentRunDetails | undefined
  >(undefined);
  const [error, setError] = useState(undefined);

  const { runID } = useParams();
  const navigate = useNavigate();

  function makeButtons() {
    const goBack = () => {
      navigate("/experimentation/list");
    };
    const goBackButton = <Button text="Back" variant="neutral" onClick={goBack} />;

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
    const destructiveButton = (
      <Button
        text={title}
        variant="destructive"
        onClick={() => {
          client
            .post("/v1/chaos/experimentation/cancelExperimentRun", { id: runID })
            .then(() => {
              setExperiment(undefined);
            })
            .catch((err: ClutchError) => {
              setError(err);
            });
        }}
      />
    );

    return [goBackButton, destructiveButton];
  }

  if (experiment === undefined && error === "") {
    client
      .post("/v1/chaos/experimentation/getExperimentRunDetails", { id: runID })
      .then(response => {
        setExperiment(response?.data?.runDetails);
      })
      .catch((err: ClutchError) => {
        setError(err);
      });
  }

  return (
    <PageLayout heading={heading} error={error}>
      <Form>
        {experiment && (
          <>
            {experiment.properties.items.map(property =>
              property.urlValue !== undefined ? (
                <Link href={property.urlValue} key={property.urlValue} textTransform="capitalize">
                  {property.label}
                </Link>
              ) : (
                <TextField
                  readOnly
                  key={property.label}
                  label={property.label}
                  defaultValue={propertyToString(property)}
                />
              )
            )}
            <TextField
              multiline
              readOnly
              label="Config"
              defaultValue={JSON.stringify(experiment.config, null, 4)}
            />
            <ButtonGroup>{makeButtons()}</ButtonGroup>
          </>
        )}
      </Form>
    </PageLayout>
  );
};

export default ViewExperimentRun;
