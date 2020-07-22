import React from "react";
import { ExpandableRow, ExpandableTable, ExpansionPanel, Status, StatusRow } from "@clutch-sh/core";
import { Grid } from "@material-ui/core";

const RatioStatus = ({ succeeded, failed, ...props }) => (
  <Grid container alignItems="center" justify="flex-end" {...props}>
    {succeeded ? (
      <Grid item>
        <Status variant="success" {...props}>
          {succeeded}
        </Status>
      </Grid>
    ) : null}
    {succeeded && failed ? <Grid item> / </Grid> : null}
    {failed ? (
      <Grid item>
        <Status variant="failure" {...props}>
          {failed}
        </Status>
      </Grid>
    ) : null}
  </Grid>
);

const clusterStatuses = data => {
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

const Clusters = ({ clusters }) => {
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
      <ExpandableTable headings={["Name", "Hosts"]}>
        {statuses.map(cluster => (
          <ExpandableRow
            key={cluster.name}
            heading={cluster.name}
            summary={
              cluster.hosts.length === 0 ? (
                <Status align="right">0</Status>
              ) : (
                <RatioStatus
                  align="right"
                  succeeded={cluster.healthyCount}
                  failed={cluster.unhealthyCount}
                />
              )
            }
          >
            {cluster.hosts.map(host => {
              const hostData = { ...host };
              const { healthy } = hostData;
              delete hostData.healthy;
              return (
                <StatusRow key={host.address} success={healthy} data={Object.values(hostData)} />
              );
            })}
          </ExpandableRow>
        ))}
      </ExpandableTable>
    </ExpansionPanel>
  );
};

export default Clusters;
