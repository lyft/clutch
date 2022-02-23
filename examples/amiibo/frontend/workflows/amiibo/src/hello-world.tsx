import type { ChangeEvent } from "react";
import React from "react";
import {
  Button,
  ButtonGroup,
  client,
  Table,
  TableRow,
  TextField,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import _ from "lodash";

import type { WorkflowProps } from ".";

const AmiiboLookup: React.FC<WizardChild> = () => {
  const { onSubmit } = useWizardContext();
  const userInput = useDataLayout("userInput");

  const onChange = (event: ChangeEvent<{ value: string }>) => {
    userInput.assign({ name: event.target.value });
  };

  return (
    <>
      <TextField onChange={onChange} onReturn={onSubmit} />
      <ButtonGroup>
        <Button text="Search" onClick={onSubmit} />
      </ButtonGroup>
    </>
  );
};

const AmiiboDetails: React.FC<WizardChild> = () => {
  const amiiboData = useDataLayout("amiiboData");
  let amiiboResults = amiiboData.displayValue();
  if (_.isEmpty(amiiboResults)) {
    amiiboResults = [];
  }

  return (
    <WizardStep error={amiiboData.error} isLoading={amiiboData.isLoading}>
      <Table columns={["Name", "Image", "Series", "Type"]}>
        {amiiboResults.map(amiibo => (
          <TableRow key={amiibo.imageUrl}>
            {amiibo.name}
            <img alt={amiibo.name} src={amiibo.imageUrl} height="75px" />
            {amiibo.amiiboSeries}
            {amiibo.type}
          </TableRow>
        ))}
      </Table>
    </WizardStep>
  );
};

const HelloWorld: React.FC<WorkflowProps> = ({ heading }) => {
  const dataLayout = {
    userInput: {},
    amiiboData: {
      deps: ["userInput"],
      hydrator: (userInput: { name: string }) => {
        return client
          .post("/v1/amiibo/getAmiibo", {
            name: userInput.name,
          })
          .then(response => {
            return response?.data?.amiibo || [];
          });
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <AmiiboLookup name="Lookup" />
      <AmiiboDetails name="Details" />
    </Wizard>
  );
};

export default HelloWorld;
