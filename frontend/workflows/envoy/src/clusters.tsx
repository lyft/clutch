import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { AccordionRow, ExpansionPanel, StatusIcon, Table, TableRow } from "@clutch-sh/core";
import { Grid } from "@material-ui/core";

interface StatusRowProps {
  success: boolean;
  data: any[];
}

const StatusRow: React.FC<StatusRowProps> = ({ success, data }) => {
  const variant = success ? "success" : "failure";
  return (
    <TableRow>
      {[...data]}
      <StatusIcon variant={variant} />
    </TableRow>
  );
};

interface RatioStatusProps {
  succeeded: boolean;
  failed: boolean;
  align?: "right" | "center";
}

const RatioStatus: React.FC<RatioStatusProps> = ({ succeeded, failed, ...props }) => (
  <Grid container alignItems="center" {...props}>
    {succeeded ? (
      <Grid item>
        <StatusIcon variant="success" {...props}>
          {succeeded}
        </StatusIcon>
      </Grid>
    ) : null}
    {succeeded && failed ? <Grid item> / </Grid> : null}
    {failed ? (
      <Grid item>
        <StatusIcon variant="failure" {...props}>
          {failed}
        </StatusIcon>
      </Grid>
    ) : null}
  </Grid>
);

const clusterStatuses = (data: IClutch.envoytriage.v1.IClusters) => {
  return data.clusterStatuses.map(clusterStatus => {
    const healthyCount = clusterStatus.hostStatuses.filter(hostStatus => hostStatus.healthy).length;
    const unhealthyCount = clusterStatus.hostStatuses.length - healthyCount;
    return {
      name: clusterStatus.name,
      healthyCount,
      unhealthyCount,
      hosts: clusterStatus.hostStatuses,
    };
  });
};

interface ClustersProps {
  clusters: IClutch.envoytriage.v1.IClusters;
}

const Clusters: React.FC<ClustersProps> = ({ clusters }) => {
  const [statuses, setStatuses] = React.useState([]);
  const [summary, setSummary] = React.useState("");

  React.useEffect(() => {
    setStatuses(clusterStatuses(clusters));
  }, [clusters]);

  React.useEffect(() => {
    const healthyHostCount = statuses
      .map(cluster => cluster.healthyCount)
      .reduce((a, b) => a + b, 0);
    const totalHostCount = statuses.map(cluster => cluster.hosts.length).reduce((a, b) => a + b, 0);
    setSummary(`(${healthyHostCount}/${totalHostCount} healthy)`);
  }, [statuses]);

  return (
    <ExpansionPanel heading="Clusters" summary={summary}>
      <Table headings={["Name", "Hosts"]}>
        {statuses.map(cluster => (
          <AccordionRow
            key={cluster.name}
            headings={[
              cluster.name,
              cluster.hosts.length === 0 ? (
                <StatusIcon>0</StatusIcon>
              ) : (
                <RatioStatus succeeded={cluster.healthyCount} failed={cluster.unhealthyCount} />
              ),
            ]}
          >
            {cluster.hosts.map(host => {
              const hostData = { ...host };
              const { healthy } = hostData;
              delete hostData.healthy;
              return (
                <StatusRow key={host.address} success={healthy} data={Object.values(hostData)} />
              );
            })}
          </AccordionRow>
        ))}
      </Table>
    </ExpansionPanel>
  );
};

export default Clusters;
