import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import {
  Alert,
  Button,
  ButtonGroup,
  CheckboxPanel,
  Chip,
  client,
  Confirmation,
  MetadataTable,
  NotePanel,
  Resolver,
  Table,
  TableRow,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import { Box } from "@mui/material";
import _ from "lodash";
import { number } from "yup";

import type { ResolverChild, TableDetailsChild, WorkflowProps } from ".";

const TableIdentifier: React.FC<ResolverChild> = ({ resolverType, notes = [] }) => {
  const { onSubmit } = useWizardContext();
  const resolvedResourceData = useDataLayout("resourceData");
  const capacityUpdates = useDataLayout("capacityUpdates");

  const onResolve = ({ results }) => {
    // Decide how to process results.
    resolvedResourceData.assign(results[0]);
    capacityUpdates.updateData(
      "gsi_map",
      _.mapValues(_.keyBy(results[0].globalSecondaryIndexes, "name"), v =>
        v.provisionedThroughput.toJSON()
      )
    );
    onSubmit();
  };

  const resolverNotes = notes.filter(note => note.location === "resolver");

  return (
    <Resolver type={resolverType} searchLimit={1} onResolve={onResolve} notes={resolverNotes} />
  );
};

const TableDetails: React.FC<TableDetailsChild> = ({ enableOverride, notes = [] }) => {
  const { onSubmit, onBack } = useWizardContext();
  const resourceData = useDataLayout("resourceData");
  const capacityUpdates = useDataLayout("capacityUpdates");
  const table = resourceData.displayValue() as IClutch.aws.dynamodb.v1.Table;

  const tableDetailsNotes = notes.filter(note => note.location === "table-details");
  const limitsNotes = notes.filter(note => note.location === "scaling-limits");
  const getChipStatus = (status: IClutch.aws.dynamodb.v1.Table.Status) => {
    switch (status.toString()) {
      case "ACTIVE":
        return "active";
      case "UPDATING":
      case "CREATING":
        return "pending";
      case "DELETING":
        return "warn";
      default:
        return "neutral";
    }
  };

  const handleTableCapacityChange = (key: string, value: number) => {
    const newTableThroughput = { ...capacityUpdates.displayValue().table_throughput, [key]: value };
    capacityUpdates.updateData("table_throughput", newTableThroughput);
  };

  const handleGsiCapacityChange = (key: string, value: number) => {
    // big hack to retrieve the capacity type (read or write?) and
    // the GSI name from a single event attribute [key]
    // where key is formatted like "read,gsi-name"
    // feature request to address this: https://github.com/lyft/clutch/issues/1739
    const [capacityType, gsiName] = key.split(",");

    const updatesList = { ...(capacityUpdates.displayValue()?.gsi_updates || {}) };
    if (!_.has(updatesList, gsiName)) {
      // copy over current throughput on first update
      updatesList[gsiName] = capacityUpdates.displayValue().gsi_map[gsiName];
      capacityUpdates.updateData("gsi_updates", updatesList);
    }
    updatesList[gsiName] = { ...updatesList[gsiName], [capacityType]: value };
    capacityUpdates.updateData("gsi_updates", updatesList);
  };

  if (!_.has(table, "globalSecondaryIndexes")) {
    table.globalSecondaryIndexes = [];
  }

  return (
    <WizardStep error={resourceData.error} isLoading={resourceData.isLoading}>
      <NotePanel notes={tableDetailsNotes} />
      <Box>
        <Table columns={["Name", "Type", "Status", "Provisioned Capacities"]}>
          <TableRow key={table.name}>
            {table.name}
            Table
            <Chip variant={getChipStatus(table.status)} label={table.status} size="small" />
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
                      .transform(value => (Number.isNaN(value) ? 0 : value))
                      .min(Number(table.provisionedThroughput.readCapacityUnits)),
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
                      .transform(value => (Number.isNaN(value) ? 0 : value))
                      .min(Number(table.provisionedThroughput.writeCapacityUnits)),
                  },
                },
              ]}
            />
          </TableRow>
          {table.globalSecondaryIndexes.map(gsi => (
            <TableRow key={gsi.name}>
              {gsi.name}
              Index
              <Chip variant={getChipStatus(gsi.status)} label={gsi.status} size="small" />
              <MetadataTable
                onUpdate={handleGsiCapacityChange}
                data={[
                  {
                    name: "read",
                    value: gsi.provisionedThroughput.readCapacityUnits,
                    input: {
                      type: "number",
                      key: `readCapacityUnits,${gsi.name}`,
                      validation: number()
                        .integer()
                        .transform(value => (Number.isNaN(value) ? 0 : value))
                        .min(Number(gsi.provisionedThroughput.readCapacityUnits)),
                    },
                  },
                  {
                    name: "write",
                    value: gsi.provisionedThroughput.writeCapacityUnits,
                    input: {
                      type: "number",
                      key: `writeCapacityUnits,${gsi.name}`,
                      validation: number()
                        .integer()
                        .transform(value => (Number.isNaN(value) ? 0 : value))
                        .min(Number(gsi.provisionedThroughput.writeCapacityUnits)),
                    },
                  },
                ]}
              />
            </TableRow>
          ))}
        </Table>
      </Box>

      {enableOverride && (
        <Box>
          {limitsNotes.length === 0 && (
            <Alert severity="warning">
              Warning: to override the DynamoDB scaling limits, check the box below. This will
              bypass the maximum limits placed on all throughput updates. Only override limits if
              safe to do so.
            </Alert>
          )}

          {limitsNotes.length > 0 && <NotePanel notes={limitsNotes} />}

          <CheckboxPanel
            onChange={state =>
              capacityUpdates.updateData("ignore_maximums", state["Override limits"])
            }
            options={{
              "Override limits": false,
            }}
          />
        </Box>
      )}

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
      <Confirmation action="Update" />
      <Table columns={["Name", "Status"]}>
        {statusList.map(s => (
          <TableRow key={s.name}>
            {s.name}
            {s.status}
          </TableRow>
        ))}
      </Table>
    </WizardStep>
  );
};

const UpdateCapacity: React.FC<WorkflowProps> = ({
  resolverType,
  notes = [],
  enableOverride,
  heading,
}) => {
  const dataLayout = {
    resourceData: {},
    capacityUpdates: {},
    updateCapacityOutput: {
      deps: ["resourceData", "capacityUpdates"],
      hydrator: (resourceData, capacityUpdates) => {
        const tableArgs = {
          tableName: resourceData.name,
          account: resourceData.account,
          region: resourceData.region,
          ignoreMaximums:
            "ignore_maximums" in capacityUpdates ? capacityUpdates.ignore_maximums : false,
        };

        let changeArgs: {};
        if (_.has(capacityUpdates, "table_throughput")) {
          const targetTableThroughput = {
            readCapacityUnits:
              capacityUpdates.table_throughput?.read ??
              resourceData.provisionedThroughput.readCapacityUnits,
            writeCapacityUnits:
              capacityUpdates.table_throughput.write ??
              resourceData.provisionedThroughput.writeCapacityUnits,
          };
          if (!_.isEqual(targetTableThroughput, resourceData.provisionedThroughput.toJSON())) {
            changeArgs = { tableThroughput: targetTableThroughput };
          }
        }
        if (_.has(capacityUpdates, "gsi_updates")) {
          const updatesFormatted = [];
          _.each(capacityUpdates.gsi_updates, (throughput, name) => {
            if (!_.isEqual(capacityUpdates.gsi_map[name], throughput)) {
              updatesFormatted.push({ name, index_throughput: throughput });
            }
          });
          changeArgs = { ...changeArgs, gsi_updates: updatesFormatted };
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
      <TableIdentifier name="Lookup" resolverType={resolverType} notes={notes} />
      <TableDetails name="Modify" enableOverride={enableOverride} notes={notes} />
      <Confirm name="Results" />
    </Wizard>
  );
};

export default UpdateCapacity;
