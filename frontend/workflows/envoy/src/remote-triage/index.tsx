import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import type { BaseWorkflowProps } from "@clutch-sh/core";
import {
  Button,
  ButtonGroup,
  client,
  Grid,
  MetadataTable,
  Paper,
  Tab,
  Tabs,
  TextField,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import styled from "@emotion/styled";

import Clusters from "./clusters";
import Listeners from "./listeners";
import Runtime from "./runtime";
import ServerInfo from "./server-info";
import Stats from "./stats";

const INCLUDE_OPTIONS = {
  clusters: true,
  configDump: true,
  listeners: true,
  runtime: true,
  stats: true,
  serverInfo: true,
};

const TriageIdentifier: React.FC<WizardChild> = () => {
  const { onSubmit } = useWizardContext();
  const resourceData = useDataLayout("resourceData");

  return (
    <>
      <TextField
        label="IP Address or Hostname"
        placeholder="127.0.0.1"
        onChange={e => resourceData.updateData("host", e.target.value)}
        onReturn={onSubmit}
      />
      <ButtonGroup>
        <Button text="Search" onClick={onSubmit} />
      </ButtonGroup>
    </>
  );
};

interface DashboardTabProps {
  serverInfo: IClutch.envoytriage.v1.IServerInfo;
  summaries?: {
    name: string;
    value: number;
  }[];
}

const SummaryCardTitle = styled.div({
  fontWeight: 600,
  fontSize: "14px",
  color: "#0D1030",
});

const SummaryCardBody = styled.div({
  fontWeight: "bold",
  fontSize: "20px",
  color: "#3548D4",
});

const DashboardTab = ({ serverInfo, summaries }: DashboardTabProps) => {
  const INFORMATION_KEYS = [
    "hot_restart_version",
    "uptime_all_epochs",
    "uptime_current_epoch",
    "version",
  ];

  const serverData = INFORMATION_KEYS.map(key => {
    return { name: key, value: serverInfo.value?.[key] };
  });

  return (
    <div>
      <Grid container direction="row" justify="space-evenly" wrap="nowrap" spacing={1}>
        <Grid item style={{ flexBasis: "60%" }}>
          <Paper>
            <SummaryCardTitle>Clusters</SummaryCardTitle>
          </Paper>
        </Grid>
        <Grid
          item
          container
          direction="column"
          justify="space-evenly"
          spacing={1}
          style={{ textAlign: "center", flexBasis: "40%" }}
        >
          {summaries.map(summary => (
            <Grid item key={summary.name}>
              <Paper>
                <SummaryCardTitle>{summary.name}</SummaryCardTitle>
                <SummaryCardBody>{summary.value}</SummaryCardBody>
              </Paper>
            </Grid>
          ))}
        </Grid>
      </Grid>
      <div style={{ padding: "16px 0" }}>
        <MetadataTable data={serverData} />
      </div>
    </div>
  );
};

const TriageDetails: React.FC<WizardChild> = () => {
  const remoteData = useDataLayout("remoteData");
  const metadata = remoteData.value.nodeMetadata as IClutch.envoytriage.v1.NodeMetadata;
  const { clusters, listeners, runtime, stats, serverInfo } =
    (remoteData.value?.output as IClutch.envoytriage.v1.Result.Output) || {};

  const summaryData = [
    { name: "Clusters", value: clusters?.clusterStatuses?.length || 0 },
    { name: "Listeners", value: listeners?.listenerStatuses?.length || 0 },
    { name: "Runtime Keys", value: runtime?.entries?.length || 0 },
    { name: "Stats", value: stats?.stats?.length || 0 },
  ];
  return (
    <WizardStep error={remoteData.error} isLoading={remoteData.isLoading}>
      <MetadataTable
        data={[
          {
            name: "Address",
            value: `${remoteData.value.address?.host}:${remoteData.value.address?.port}`,
          },
          { name: "Service Node", value: metadata?.serviceNode },
          { name: "Service Zone", value: metadata?.serviceZone },
          { name: "Service Cluster", value: metadata?.serviceCluster },
          // { name: "Version", value: metadata?.version },
        ]}
      />
      <Tabs>
        <Tab label="Dashboard">
          <DashboardTab serverInfo={serverInfo} summaries={summaryData} />
        </Tab>
        <Tab label="Clusters">
          <Clusters clusters={clusters} />
        </Tab>
        <Tab label="Listeners">
          <Listeners listeners={listeners} />
        </Tab>
        <Tab label="Runtime">
          <Runtime runtime={runtime} />
        </Tab>
        <Tab label="Stats">
          <Stats stats={stats} />
        </Tab>
        <Tab label="Server Info">
          <ServerInfo info={serverInfo} />
        </Tab>
      </Tabs>
    </WizardStep>
  );
};

const RemoteTriage: React.FC<BaseWorkflowProps> = ({ heading }) => {
  const dataLayout = {
    resourceData: {},
    remoteData: {
      deps: ["resourceData"],
      cache: false,
      hydrator: (resourceData: { host: string }) => {
        return client.post("/v1/envoytriage/read", {
          operations: [
            {
              address: {
                host: resourceData.host,
              },
              include: INCLUDE_OPTIONS,
            },
          ],
        } as IClutch.envoytriage.v1.ReadRequest);
      },
      transformResponse: response => response.data.results?.[0],
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <TriageIdentifier name="Lookup" />
      <TriageDetails name="Details" />
    </Wizard>
  );
};

export default RemoteTriage;
