import React from "react";
import { useForm } from "react-hook-form";
import type { clutch as IClutch } from "@clutch-sh/api";
import {
  Button,
  ButtonGroup,
  client,
  Confirmation,
  MetadataTable,
  Resolver,
  Select,
  TextField,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import styled from "@emotion/styled";
import { Grid } from "@material-ui/core";

import type { ResolverChild, WorkflowProps } from "./index";

const Form = styled.form({
  alignItems: "center",
  display: "flex",
  flexDirection: "column",
  justifyItems: "space-evenly",
  "> *": {
    padding: "8px 0",
  },
});

const TableIdentifier: React.FC<ResolverChild> = ({ resolverType }) => {
  const { onSubmit } = useWizardContext();
  const resolvedResourceData = useDataLayout("resourceData");

  const onResolve = ({ results }) => {
    // Decide how to process results.
    resolvedResourceData.assign(results[0]);
    onSubmit();
  };

  return <Resolver type={resolverType} searchLimit={1} onResolve={onResolve} />;
};

const TableDetails: React.FC<WizardChild> = () => {
  const { handleSubmit } = useForm({
    mode: "onChange",
  });
  const { onSubmit, onBack } = useWizardContext();
  const resourceData = useDataLayout("resourceData");
  const table = resourceData.displayValue() as IClutch.aws.dynamodb.v1.Table;
  const handleTargetCapacityChange = (event) => {
    resourceData.updateData(event.target.name, event.target.value);
  };


  return (
    <WizardStep error={resourceData.error} isLoading={resourceData.isLoading}>
      <Form onSubmit={handleSubmit(onSubmit)}>
        <TextField readOnly label="TableName" name="TableName" value={table.name} />
        <TextField readOnly label="Region" name="region" value={table.region} />
        <Grid container alignItems="stretch" wrap="nowrap">
          <Grid item style={{ flexBasis: "50%", paddingRight: "8px" }}>
            <TextField
              readOnly
              label="Current RCU"
              name="currentReadCapacityUnits"
              value={table.provisionedThroughput.readCapacityUnits}
              disabled
            />
            <TextField
              readOnly
              label="Current WCU"
              name="currentWriteCapacityUnits"
              value={table.provisionedThroughput.writeCapacityUnits}
              disabled
            />
          </Grid>
          <Grid item style={{ flexBasis: "50%", paddingLeft: "8px" }}>
            <Select
              label="TargetReadCapacityUnits"
              name="targetReadCapacityUnits"
              onChange={handleTargetCapacityChange}
            />
            <Select
              label="TargetWriteCapacityUnits"
              name="targetWriteCapacityUnits"
              onChange={handleTargetCapacityChange}
            />
          </Grid>
        </Grid>
      </Form>

      <ButtonGroup>
        <Button text="Back" onClick={() => onBack()} />
        <Button text="Update" variant="destructive" onClick={onSubmit} />
      </ButtonGroup>
    </WizardStep>
  );
};

const Confirm: React.FC<WizardChild> = () => {
  const updateTableData = useDataLayout("tableData");
  const tableData = useDataLayout("resourceData").value;
  return (
    <WizardStep error={updateTableData.error} isLoading={updateTableData.isLoading}>
      <Confirmation action="Update" />
      <MetadataTable
        data={[
          { name: "Table Name", value: tableData.tableName },
          { name: "Region", value: tableData.region },
          { name: "Current Read Capacity", value: tableData.currentReadCapacityUnits },
          { name: "Target Read Capacity", value: tableData.targetReadCapacityUnits },
        ]}
      />
      <MetadataTable
        data={[
          { name: "Table Name", value: tableData.tableName },
          { name: "Region", value: tableData.region },
          { name: "Current Write Capacity", value: tableData.currentWriteCapacityUnits },
          { name: "Target Write Capacity", value: tableData.targetWriteCapacityUnits },
        ]}
      />
    </WizardStep>
  );
};

const UpdateTableCapacity: React.FC<WorkflowProps> = ({ heading, resolverType }) => {
  const dataLayout = {
    resourceData: {},
    streamData: {
      deps: ["resourceData"],
      hydrator: (resourceData: {
        tableName: string;
        region: string;
        targetWriteCapacityUnits: number;
        targetReadCapacityUnits: number;
      }) => {
        return client.post("/v1/aws/dynamodb/updateTableCapacity", {
          stream_name: resourceData.tableName,
          region: resourceData.region,
          target_table_rcu: resourceData.targetReadCapacityUnits,
          target_table_wcu: resourceData.targetWriteCapacityUnits,
        });
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <TableIdentifier name="Lookup" resolverType={resolverType} />
      <TableDetails name="Modify" />
      <Confirm name="Confirmation" />
    </Wizard>
  );
};

export default UpdateTableCapacity;
