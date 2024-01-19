import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { AccordionRow, StatusIcon, styled, Table, TableRow, useTheme } from "@clutch-sh/core";
import type { Theme } from "@mui/material";
import _ from "lodash";

const BarContainer = styled("rect")<{ $fill: string; $width: string }>(
  {
    height: "12px",
  },
  props => ({ theme }: { theme: Theme }) => ({
    width: props.$width,
    fill: props.$fill,
    strokeWidth: props.$fill === "transparent" ? "1px" : "0",
    stroke: theme.palette.primary[400],
  })
);

const Bar = ({ fill, width }) => (
  <svg width={width} height="12px">
    <BarContainer $fill={fill} $width={width} />
  </svg>
);

interface RatioStatusProps {
  succeeded: number;
  failed: number;
}

const RatioStatus: React.FC<RatioStatusProps> = ({ succeeded, failed }) => {
  const theme = useTheme();
  const total = succeeded + failed;
  return (
    <>
      {succeeded !== 0 && (
        <Bar fill={theme.palette.success[300]} width={`${(succeeded / total) * 100}px`} />
      )}
      {failed !== 0 && (
        <Bar fill={theme.palette.error[300]} width={`${(failed / total) * 100}px`} />
      )}
    </>
  );
};

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

interface StatusRowProps {
  success: boolean;
  data: any[];
}

export const StatusRow = ({ success, data }: StatusRowProps) => {
  const displayData = [...data];
  const headerValue = displayData.shift();
  const variant = success ? "success" : "failure";
  return (
    <TableRow>
      <div style={{ textAlign: "center" }}>{headerValue}</div>
      <StatusIcon align="left" variant={variant} />
    </TableRow>
  );
};

interface ClustersProps {
  clusters: IClutch.envoytriage.v1.IClusters;
}

const Clusters: React.FC<ClustersProps> = ({ clusters }) => {
  const [statuses, setStatuses] = React.useState([]);

  React.useEffect(() => {
    setStatuses(clusterStatuses(clusters));
  }, [clusters]);

  return (
    <div style={{ height: "400px", display: "flex" }}>
      <Table stickyHeader columns={["Hosts", "Status"]}>
        {_.sortBy(statuses, ["name"]).map(cluster => (
          <AccordionRow
            key={cluster.name}
            columns={[
              cluster.name,
              cluster.hosts.length === 0 ? (
                <Bar fill="transparent" width="100px" />
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
    </div>
  );
};

export default Clusters;
