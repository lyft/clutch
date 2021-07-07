import React from "react";
import TimeAgo from "react-timeago";
import type { clutch as IClutch } from "@clutch-sh/api";
import {
  Chip,
  Table,
  TableRow,
  TableRowAction,
  TableRowActions,
  useNavigate,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import styled from "@emotion/styled";
import DeleteIcon from "@material-ui/icons/Delete";
import _ from "lodash";

const PodsContainer = styled.div({
  display: "flex",
  maxHeight: "50vh",
});

const getReadyCountString = containers => {
  const readyCount = containers.filter(cont => cont.ready).length;
  return `${readyCount.toString()}/${containers?.length?.toString()}`;
};

const getRestartCountString = containers => {
  const restartCount = containers.reduce((a, b) => a + b.restartCount, 0);
  return restartCount.toString();
};

const getPodChipState = (status: string) => {
  switch (true) {
    case status.includes("Error") ||
      status.includes("BackOff") ||
      status.includes("Evicted") ||
      status.includes("Terminating"):
      return "error";
    case status.includes("Running"):
      return "active";
    case status.includes("Init") || status.includes("Initializing") || status.includes("Pending"):
      return "pending";
    case status.includes("Completed"):
      return "success";
    default:
      return "neutral";
  }
};

const convertTime = timeStampMillis => {
  return parseInt(timeStampMillis, 10);
};

const timeFormatter = (value, unit, suffix) => {
  if (suffix === "ago") {
    return `${value}${unit.charAt(0)}`;
  }
  return `${value}${unit.charAt(0)} ${suffix}`;
};

const PodTable = () => {
  const podListData = useDataLayout("podListData", { hydrate: false });
  const pods = podListData.displayValue()?.pods as IClutch.k8s.v1.Pod[];
  const navigate = useNavigate();

  return (
    <PodsContainer>
      <Table
        stickyHeader
        actionsColumn
        headings={[
          "Name",
          "Cluster",
          "Containers Ready",
          "Restart",
          "Node IP",
          "Pod IP",
          "State",
          "SHA",
          "Age",
        ]}
      >
        {_.sortBy(pods, [
          o => {
            return o.name;
          },
        ]).map(pod => (
          <TableRow key={pod.name} defaultCellValue="nil">
            {pod.name}
            {pod.cluster}
            {getReadyCountString(pod.containers)}
            {getRestartCountString(pod.containers)}
            {pod.nodeIp}
            {pod.podIp}
            <Chip variant={getPodChipState(pod.status)} label={pod.status} />
            {pod.labels?.version}
            <TimeAgo date={convertTime(pod.startTimeMillis)} formatter={timeFormatter} />
            <TableRowActions>
              <TableRowAction
                icon={<DeleteIcon />}
                onClick={() =>
                  navigate(`/k8s/pod/delete?q=${pod.cluster}/${pod.namespace}/${pod.name}`, {
                    origin: true,
                  })
                }
              >
                Delete
              </TableRowAction>
            </TableRowActions>
          </TableRow>
        ))}
      </Table>
    </PodsContainer>
  );
};

export { PodTable as default, timeFormatter, convertTime };
