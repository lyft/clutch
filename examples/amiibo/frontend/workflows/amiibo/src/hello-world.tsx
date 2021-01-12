import React, { ChangeEvent } from "react";
import _ from "lodash";
import {
  Button,
  ButtonGroup,
  client,
  useWizardContext,
  Table,
  TableRow,
  TextField,
} from "@clutch-sh/core";

import { useDataLayout } from "@clutch-sh/data-layout";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import type { WizardChild } from "@clutch-sh/wizard";

import type { WorkflowProps } from ".";

const AmiiboLookup: React.FC<WizardChild> = () => {
  const { onSubmit } = useWizardContext();
  const userInput = useDataLayout("userInput");

  const onChange = ((event: ChangeEvent<{value: string}>) => {
    userInput.assign({name: event.target.value});
  });

  return (
    <>
      <TextField onChange={onChange} onReturn={onSubmit}/>
      <ButtonGroup>
        <Button text="Search" onClick={onSubmit}/>
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
      <Table headings={["Name", "Image", "Series", "Type"]}>
        {amiiboResults.map((amiibo, index: number) => (
          <TableRow key={index}>
            {amiibo.name}
            <img src={amiibo.imageUrl} height="75px"/>
            {amiibo.amiiboSeries}
            {amiibo.type}
          </TableRow>
        ))}
      </Table>
    </WizardStep>
  );
};

const Amiibo: React.FC<WorkflowProps> = ({ heading }) => {
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

export default Amiibo;
