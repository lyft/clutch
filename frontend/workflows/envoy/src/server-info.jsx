import React from "react";
import { ExpansionPanel, MetadataTable } from "@clutch-sh/core";
import { TableCell, TableRow } from "@material-ui/core";

const ServerInfo = ({ info }) => {
  const serverInfo = { ...info.value };
  const cliOptions = serverInfo.command_line_options;
  delete serverInfo.command_line_options;

  const status = `(${serverInfo.state.toLowerCase()})`;
  delete serverInfo.state;

  const serverData = Object.keys(serverInfo).map(key => {
    return { name: key, value: serverInfo[key] };
  });
  serverData.push({ name: "Command Line Options", value: "" });
  const cliOptionMetadata = Object.keys(cliOptions).map(key => {
    return { name: key, value: cliOptions[key] };
  });
  const midPoint = Math.floor(cliOptionMetadata.length / 2);
  const variant = "small";
  return (
    <ExpansionPanel heading="Server Info" summary={status}>
      <MetadataTable data={serverData} variant="small">
        <TableRow>
          <TableCell size={variant}>
            <MetadataTable data={cliOptionMetadata.slice(0, midPoint)} variant={variant} />
          </TableCell>
          <TableCell size={variant}>
            <MetadataTable data={cliOptionMetadata.slice(midPoint)} variant={variant} />
          </TableCell>
        </TableRow>
      </MetadataTable>
    </ExpansionPanel>
  );
};

export default ServerInfo;
