import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { MetadataTable, styled } from "@clutch-sh/core";
import type { Theme } from "@mui/material";

const Container = styled("div")(({ theme }: { theme: Theme }) => ({
  "> *": {
    padding: theme.spacing("sm", "none"),
  },
}));

const Title = styled("div")(({ theme }: { theme: Theme }) => ({
  fontWeight: "bold",
  fontSize: "20px",
  color: theme.palette.secondary[900],
  textTransform: "capitalize",
}));

interface ServerInformation {
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
  const status = `${rawServerInfo.state.toLowerCase()}`;

  const serverData = Object.keys(information).map(key => {
    return { name: key, value: information[key] };
  });
  const cliOptionMetadata = Object.keys(cliOptions).map(key => {
    return { name: key, value: cliOptions[key] };
  });
  return (
    <Container>
      <Title>{status}</Title>
      <Title>This is a test</Title>
      <MetadataTable data={serverData} />
      <Title>Command Line Options</Title>
      <MetadataTable maxHeight="400px" data={cliOptionMetadata} />
    </Container>
  );
};

export default ServerInfo;
