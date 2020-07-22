import React from "react";
import { Button, client, TextField, useWizardContext } from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import { Typography } from "@material-ui/core";

const EchoInput = () => {
  const { onSubmit } = useWizardContext();
  const echoInput = useDataLayout("echoInput");

  return (
    <WizardStep>
      <TextField
        InputLabelProps={ { color: "secondary" } }
        label="Enter in text to echo back"
        name="echoData"
        required
        onChange={e => echoInput.assign({ [e.target.name]: e.target.value })}
      />
      <Button text="Echo" destructive onClick={onSubmit} />
    </WizardStep>
  );
};

const EchoOutput = () => {
  const echoData = useDataLayout("echoOutput");

  return (
    <WizardStep error={echoData.error} isLoading={echoData.isLoading}>
      <Typography>{echoData.message}</Typography>
    </WizardStep>
  );
};

const EchoWorkflow = ({ heading }) => {
  const dataLayout = {
    echoInput: {},
    echoOutput: {
      deps: ["echoInput"],
      hydrator: echoInput => {
        return client
          .post("/v1/echo/sayHello", {
            name: echoInput.name,
          })
          .then(resp => {
            return resp;
          });
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <EchoInput name="Lookup" />
      <EchoOutput name="Echo" heading="Instance Details" />
    </Wizard>
  );
};

export default EchoWorkflow;
