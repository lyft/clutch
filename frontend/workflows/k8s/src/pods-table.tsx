import React from "react";
import { useNavigate } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Table, TableRow, TableRowAction, TableRowActions } from "@clutch-sh/core";
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

const PodTable = () => {
  const podListData = useDataLayout("podListData");
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
            {pod.status}
            {pod.labels?.version}
            <TableRowActions>
              <TableRowAction
                icon={<DeleteIcon />}
                onClick={() =>
                  navigate(`/k8s/pod/delete?q=${pod.cluster}/${pod.namespace}/${pod.name}`)
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

export default PodTable;
