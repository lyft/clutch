import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import {
  Button,
  ButtonGroup,
  CheckboxPanel,
  client,
  MetadataTable,
  Resolver,
  Table,
  TableRow,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import styled from "@emotion/styled";
import _ from "lodash";
import { number } from "yup";

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

  const handleTableCapacityChange = (key: string, value: string) => {
    const newTableCapacity = { ...capacityUpdates.displayValue().table_throughput };
    newTableCapacity[key] = value;
    capacityUpdates.updateData("table_throughput", newTableCapacity);
  };

  const handleGsiCapacityChange = (key: string, value: string) => {
    // big hack to retrieve the capacity type (read or write?) and
    // the GSI name from a single event attribute
    const keys = key.split(",");
    const capacityType = keys[0];
    const gsiName = keys[1];

    const gsiList = capacityUpdates.displayValue().gsi_updates
      ? [...capacityUpdates.displayValue().gsi_updates]
      : [];
    const idx = gsiList.findIndex(
      (gsi: { name: string; indexThroughput: {} }) => gsi.name === gsiName
    );
    if (idx > -1) {
      // gsi already in the edits list
      const gsi = gsiList[idx];
      gsi.index_throughput = { ...gsi.index_throughput, [capacityType]: value };
      capacityUpdates.updateData("gsi_updates", gsiList);
    } else {
      // copy over current capacities
      const curr = table.globalSecondaryIndexes.find(gsi => gsi.name === gsiName);
      const newGsi = {
        name: gsiName,
        index_throughput: {
          read_capacity_units: curr.provisionedThroughput.readCapacityUnits,
          write_capacity_units: curr.provisionedThroughput.writeCapacityUnits,
        },
      };
      newGsi.index_throughput = { ...newGsi.index_throughput, [capacityType]: value };
      gsiList.push(newGsi);
      capacityUpdates.updateData("gsi_updates", gsiList);
    }
  };

  return (
    <WizardStep error={resourceData.error} isLoading={resourceData.isLoading}>
      <Container>
        <Table columns={["Name", "Type", "Status", "Provisioned Capacities"]}>
          <TableRow key={table.name}>
            {table.name}
            Table
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
                      .integer("must be a number")
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
                },
              ]}
            />
          </TableRow>
          {table.globalSecondaryIndexes.map(gsi => (
            <TableRow key={gsi.name}>
              {gsi.name}
              Index
              {gsi.status}
              <MetadataTable
                onUpdate={handleGsiCapacityChange}
                data={[
                  {
                    name: "read",
                    value: gsi.provisionedThroughput.readCapacityUnits,
                    input: {
                      type: "number",
                      key: ["read_capacity_units", gsi.name],
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
                      key: ["write_capacity_units", gsi.name],
                      validation: number()
                        .integer()
                        .min(gsi.provisionedThroughput.writeCapacityUnits),
                    },
                  },
                ]}
              />
            </TableRow>
          ))}
        </Table>
      </Container>

      {/* TO DO: conditionally render the override checkbox depending on workflow config prop */}
      <CheckboxPanel
        header="To override the safety limits for scaling, check the box below."
        onChange={state => capacityUpdates.updateData("ignore_maximums", state)}
        options={{
          "Override maximum limits": false,
        }}
      />

      <ButtonGroup>
        <Button text="Back" variant="neutral" onClick={() => onBack()} />
        <Button text="Update" onClick={onSubmit} />
      </ButtonGroup>
    </WizardStep>
  );
};

const Confirm: React.FC<WizardChild> = () => {
  const updateCapacityOutput = useDataLayout("updateCapacityOutput");
  const scalingResults = updateCapacityOutput.displayValue()?.data?.table;
  let statusList = [];
  if (!_.isEmpty(scalingResults)) {
    statusList.push(scalingResults);
    statusList = statusList.concat(scalingResults.globalSecondaryIndexes);
  }

  return (
    <WizardStep error={updateCapacityOutput.error} isLoading={updateCapacityOutput.isLoading}>
      <Table columns={["Name", "Status"]}>
        {statusList.map((s, index: number) => (
          <TableRow key={s.name + index}>
            {s.name}
            {s.status}
          </TableRow>
        ))}
      </Table>
    </WizardStep>
  );
};

const UpdateCapacity: React.FC<WorkflowProps> = ({ resolverType }) => {
  const dataLayout = {
    resourceData: {},
    capacityUpdates: {},
    updateCapacityOutput: {
      deps: ["resourceData", "capacityUpdates"],
      hydrator: (resourceData, capacityUpdates) => {
        const tableArgs = {
          table_name: resourceData.name,
          region: resourceData.region,
          // ignore_maximums: capacityUpdates?.ignore_maximums? capacityUpdates.ignore_maximums : false,
        };

        let changeArgs: {};
        if (capacityUpdates.table_throughput) {
          changeArgs = {
            ...changeArgs,
            table_throughput: {
              read_capacity_units: capacityUpdates.table_throughput.read
                ? capacityUpdates.table_throughput.read
                : resourceData.provisionedThroughput.readCapacityUnits,
              write_capacity_units: capacityUpdates.table_throughput.write
                ? capacityUpdates.table_throughput.write
                : resourceData.provisionedThroughput.writeCapacityUnits,
            },
          };
        }
        if (Array.isArray(capacityUpdates.gsi_updates)) {
          changeArgs = { ...changeArgs, gsi_updates: [...capacityUpdates.gsi_updates] };
        }
        return client
          .post("/v1/aws/dynamodb/updateCapacity", { ...tableArgs, ...changeArgs })
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
