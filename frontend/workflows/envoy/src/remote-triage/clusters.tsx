import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { AccordionRow, StatusIcon, Table, TableRow } from "@clutch-sh/core";
import styled from "@emotion/styled";
import _ from "lodash";

const BarContainer = styled.rect(
  {
    height: "12px",
  },
  props => ({
    width: props.width,
    fill: props.fill,
    strokeWidth: props.fill === "transparent" ? "1px" : "0",
    stroke: "#C2C8F2",
  })
);

const Bar = ({ fill, width }) => (
  <svg width={width} height="12px">
    <BarContainer fill={fill} width={width} />
  </svg>
);

interface RatioStatusProps {
  succeeded: number;
  failed: number;
}

const RatioStatus: React.FC<RatioStatusProps> = ({ succeeded, failed }) => {
  const total = succeeded + failed;
  return (
    <>
      {succeeded !== 0 && <Bar fill="#69F0AE" width={`${(succeeded / total) * 100}px`} />}
      {failed !== 0 && <Bar fill="#FF8A80" width={`${(failed / total) * 100}px`} />}
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
      <StatusIcon align="center" variant={variant} />
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
      <Table stickyHeader headings={["Hosts", "Status"]}>
        {_.sortBy(statuses, ["name"]).map(cluster => (
          <AccordionRow
            key={cluster.name}
            headings={[
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
