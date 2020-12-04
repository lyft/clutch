import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { ExpansionPanel, MetadataTable } from "@clutch-sh/core";
import { TableCell, TableRow } from "@material-ui/core";

interface ServerInformation {
  // eslint-disable-next-line camelcase
  command_line_options: {
    [key: string]: string;
  };
  state: string;
}

const INFORMATION_KEYS = [
  "hot_restart_version",
  "uptime_all_epochs",
  "uptime_current_epoch",
  "version",
];

const ServerInfo: React.FC<{ info: IClutch.envoytriage.v1.IServerInfo }> = ({ info }) => {
  const rawServerInfo = { ...info.value } as ServerInformation;
  const information = INFORMATION_KEYS.reduce((filteredInfo, key) => {
    const localInfo = filteredInfo;
    localInfo[key] = rawServerInfo[key];
    return localInfo;
  }, {});
  const cliOptions = rawServerInfo.command_line_options;
  const status = `(${rawServerInfo.state.toLowerCase()})`;

  const serverData = Object.keys(information).map(key => {
    return { name: key, value: information[key] };
  });
  serverData.push({ name: "Command Line Options", value: "" });
  const cliOptionMetadata = Object.keys(cliOptions).map(key => {
    return { name: key, value: cliOptions[key] };
  });
  const midPoint = Math.floor(cliOptionMetadata.length / 2);
  const variant = "small";
  return (
    <ExpansionPanel heading="Server Info" summary={status}>
      <MetadataTable data={serverData}>
        <TableRow>
          <TableCell>
            <MetadataTable data={cliOptionMetadata.slice(0, midPoint)} />
          </TableCell>
          <TableCell>
            <MetadataTable data={cliOptionMetadata.slice(midPoint)} />
          </TableCell>
        </TableRow>
      </MetadataTable>
    </ExpansionPanel>
  );
};

export default ServerInfo;
