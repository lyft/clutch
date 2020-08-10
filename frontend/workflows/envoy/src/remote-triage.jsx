import React from "react";
import { Button, client, MetadataTable, TextField, useWizardContext } from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import { Typography } from "@material-ui/core";

import Clusters from "./clusters";
import Listeners from "./listeners";
import Runtime from "./runtime";
import ServerInfo from "./server-info";
import Stats from "./stats";

const INCLUDE_OPTIONS = {
  Clusters: "clusters",
  "Config Dump": "configDump",
  Listeners: "listeners",
  Runtime: "runtime",
  Stats: "stats",
  "Server Info": "serverInfo",
};

const TriageIdentifier = () => {
  const { onSubmit } = useWizardContext();
  const resourceData = useDataLayout("resourceData");

  return (
    <WizardStep>
      <TextField
        label="IP Address or Hostname"
        placeholder="127.0.0.1"
        onChange={e => resourceData.updateData("host", e.target.value)}
        onReturn={onSubmit}
      />
      <Button text="Search" onClick={onSubmit} />
    </WizardStep>
  );
};

const TriageDetails = () => {
  const remoteData = useDataLayout("remoteData");
  const metadata = remoteData.value.nodeMetadata;
  const details = remoteData.value.output;
  return (
    <WizardStep error={remoteData.error} isLoading={remoteData.isLoading}>
      <Typography variant="h6">
        <strong>
          {remoteData.value.address?.host}:{remoteData.value.address?.port}
        </strong>
      </Typography>
      <MetadataTable
        data={[
          { name: "Service Node", value: metadata?.serviceNode },
          { name: "Service Zone", value: metadata?.serviceZone },
          { name: "Service Cluster", value: metadata?.serviceCluster },
          { name: "Version", value: metadata?.version },
        ]}
      />
      {details?.clusters && <Clusters clusters={details.clusters} />}
      {details?.listeners && <Listeners listeners={details.listeners} />}
      {details?.runtime && <Runtime runtime={details.runtime} />}
      {details?.stats && <Stats stats={details.stats} />}
      {details?.serverInfo && <ServerInfo info={details.serverInfo} />}
    </WizardStep>
  );
};

const RemoteTriage = ({ heading, options }) => {
  const includeOptions = {};
  Object.values(options || INCLUDE_OPTIONS).forEach(option => {
    includeOptions[option] = true;
  });
  const dataLayout = {
    resourceData: {},
    remoteData: {
      deps: ["resourceData"],
      cache: false,
      hydrator: resourceData => {
        return client.post("/v1/envoytriage/read", {
          operations: [
            {
              address: {
                host: resourceData.host,
              },
              include: includeOptions,
            },
          ],
        });
      },
      transformResponse: response => response.data.results?.[0],
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading} maxWidth={false}>
      <TriageIdentifier name="Lookup" options={options} />
      <TriageDetails name="Details" heading="Details" />
    </Wizard>
  );
};

export default RemoteTriage;
