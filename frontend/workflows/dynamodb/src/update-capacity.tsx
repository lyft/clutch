import React, { useState } from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import {
  Button,
  ButtonGroup,
  client,
  Confirmation,
  MetadataTable,
  Resolver,
  Table,
  TableRow,
  TextField,
  useWizardContext,
} from "@clutch-sh/core";
import styled from "@emotion/styled";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import { Checkbox, FormControlLabel, Grid } from "@material-ui/core";
import _ from "lodash";
import { number, ref } from "yup";
import type Reference from "yup/lib/Reference";

import type { ResolverChild, WorkflowProps } from "./index";

const Container = styled.div({
  display: "flex",
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
  const { onSubmit, onBack } = useWizardContext();
  const resourceData = useDataLayout("resourceData");
  const capacityUpdates = useDataLayout("capacityUpdates");
  const table = resourceData.displayValue() as IClutch.aws.dynamodb.v1.Table;

  // const [overrideToggle, setOverrideToggle] = useState(false) ** for later work

  const handleTableCapacityChange = (key: string, value: string) => {
    const newTableCapacity = {...capacityUpdates.displayValue().table_throughput};
    console.log(key)
    newTableCapacity[key] = value;
    capacityUpdates.updateData("table_throughput", newTableCapacity)
    console.log(capacityUpdates.displayValue())
  };

  const handleGsiCapacityChange = (key: string, value: string) => {
    console.log(key)
    console.log(key[0])
    // const gsiList = capacityUpdates.displayValue().gsi_updates? {...capacityUpdates.displayValue().gsi_updates} : {};
    // if (event.target.name in gsiList) {
    //   gsiList[event.target.name] = {...gsiList[event.target.name], [event.target.name]: event.target.value}
    // } else {
    //   gsiList[event.target.name] = {[event.target.name]: event.target.value}
    // }
    // capacityUpdates.updateData("gsi_updates", gsiList)
  };

  // FOR LATER WORK
  // const handleOverrideToggleChange = e => {
  //   setOverrideToggle(e.target.checked);
  //   resourceData.updateData(e.target.name, e.target.checked);
  // };

  


  return (
    <WizardStep error={resourceData.error} isLoading={resourceData.isLoading}>
      <Container>
        <Table
          columns={[
            'Name',
            'Type',
            'Status',
            'Provisioned Capacities',
          ]}
        >
          <TableRow key={table.name}>
            {table.name}
            {"Table"}
            {table.status}
            <MetadataTable
              onUpdate={handleTableCapacityChange}
              data={[
                {
                  name: "read",
                  value: table.provisionedThroughput.readCapacityUnits,
                  input: {
                    type: "number",
                    key: "read",
                    validation: number()
                    .integer()
                    .min(table.provisionedThroughput.readCapacityUnits),
                  },
                },
                {
                  name: "write",
                  value: table.provisionedThroughput.writeCapacityUnits,
                  input: {
                    type: "number",
                    key: "write",
                    validation: number()
                      .integer()
                      .min(table.provisionedThroughput.writeCapacityUnits),
                  },
                }]}>
              </MetadataTable>
          </TableRow>
          {table.globalSecondaryIndexes.map(gsi => (
          <TableRow key={gsi.name}>
            {gsi.name}
            {"Index"}
            {gsi.status}
            <MetadataTable
              onUpdate={handleGsiCapacityChange}
              data={[
                {
                  name: "read",
                  value: gsi.provisionedThroughput.readCapacityUnits,
                  input: {
                    type: "number",
                    key: "r" + gsi.name,
                    validation: number()
                    .integer()
                    .min(gsi.provisionedThroughput.readCapacityUnits),
                  },
                },
                {
                  name: "write",
                  value: gsi.provisionedThroughput.writeCapacityUnits,
                  input: {
                    type: "number",
                    key: "w" + gsi.name,
                    validation: number()
                      .integer()
                      .min(gsi.provisionedThroughput.writeCapacityUnits),
                  },
                }]}>
              </MetadataTable>
          </TableRow>
      ))}
        
        </Table>

      </Container>


      <ButtonGroup>
        <Button text="Back" onClick={() => onBack()} />
        <Button text="Update" variant="destructive" onClick={onSubmit} />
      </ButtonGroup>
    </WizardStep>
  );
};

const Confirm: React.FC<WizardChild> = () => {
  const updateCapacityOutput = useDataLayout("updateCapacityOutput");
  console.log(updateCapacityOutput.displayValue())
  let scalingResults = updateCapacityOutput.displayValue()?.data?.table
  let statusList = [];
  if (!_.isEmpty(scalingResults)) {
    statusList.push(scalingResults)
    statusList = statusList.concat(scalingResults.globalSecondaryIndexes)
  } 

  return (
    <WizardStep error={updateCapacityOutput.error} isLoading={updateCapacityOutput.isLoading}>
      <Table columns={["Name", "Status"]}>
      {statusList.map((s, index: number) => (
          <TableRow key={index}>
            {s.name}
            {s.status}
          </TableRow>
      ))}
      </Table>
    </WizardStep>
  );
};

const UpdateCapacity: React.FC<WorkflowProps> = ({resolverType}) => {
  const dataLayout = {
    resourceData: {},
    capacityUpdates: {},
    updateCapacityOutput: {
      deps: ["resourceData", "capacityUpdates"],
      hydrator: (resourceData, capacityUpdates)=> {
        let tableArgs = {
          table_name: resourceData.name,
          region: resourceData.region,
          // ignore_maximums: capacityUpdates?.ignore_maximums? true : false,
        }

        let changeArgs: {}
        if (capacityUpdates.table_throughput) {
          changeArgs = {...changeArgs, table_throughput: {
            read_capacity_units: capacityUpdates.table_throughput["read"]? capacityUpdates.table_throughput["read"] : resourceData.provisionedThroughput.readCapacityUnits,
            write_capacity_units: capacityUpdates.table_throughput["write"]? capacityUpdates.table_throughput["write"] : resourceData.provisionedThroughput.writeCapacityUnits,
          }}
        }
        if (capacityUpdates.gsi_updates) {
          const gsi_list = []
          changeArgs = {...changeArgs, gsi_updates: gsi_list}
        }
        return client
          .post("/v1/aws/dynamodb/updateCapacity", {...tableArgs, ...changeArgs})
          .then(resp => {
            return resp;
          });
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout}>
      <TableIdentifier name="Lookup" resolverType={resolverType} />
      <TableDetails name="Modify" />
      <Confirm name="Results" />
    </Wizard>
  );
};

export default UpdateCapacity;