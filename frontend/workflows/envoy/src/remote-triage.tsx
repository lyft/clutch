import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import type { BaseWorkflowProps } from "@clutch-sh/core";
import { Button, client, MetadataTable, TextField, useWizardContext } from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import { Typography } from "@material-ui/core";

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
      <Button text="Search" onClick={onSubmit} />
    </>
  );
};

const TriageDetails: React.FC<WizardChild> = () => {
  const remoteData = useDataLayout("remoteData");
  const metadata = remoteData.value.nodeMetadata as IClutch.envoytriage.v1.NodeMetadata;
  const details = remoteData.value.output as IClutch.envoytriage.v1.Result.Output;
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
    <Wizard dataLayout={dataLayout} heading={heading} maxWidth={false}>
      <TriageIdentifier name="Lookup" />
      <TriageDetails name="Details" />
    </Wizard>
  );
};

export default RemoteTriage;
