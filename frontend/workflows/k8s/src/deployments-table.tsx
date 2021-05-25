import React from "react";
import { useNavigate } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Table, TableRow } from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import styled from "@emotion/styled";
import _ from "lodash";

const DeploymentsContainer = styled.div({
  display: "flex",
  maxHeight: "50vh",
});

const getReadyCountString = containers => {
  const readyCount = containers.filter(cont => cont.ready).length;
  return `${readyCount.toString()}/${containers?.length?.toString()}`;
};

const DeploymentTable = () => {
  const deploymentListData = useDataLayout("deploymentListData", { hydrate: false });
  const deployments = deploymentListData.displayValue()?.deployments as IClutch.k8s.v1.Deployment[];
  const navigate = useNavigate();

  return (
    <DeploymentsContainer>
      <Table
        stickyHeader
        actionsColumn
        headings={[
          "Name",
          "Cluster",
          "Replicas Ready",
          "Replicas Available",
          "Replicas Up-To-Date",
        ]}
      >
        {_.sortBy(deployments, [
          o => {
            return o.name;
          },
        ]).map(deployment => (
          <TableRow key={deployment.name} defaultCellValue="nil">
            {deployment.name}
            {deployment.cluster}
            {deployment.deploymentStatus?.readyReplicas}
            {deployment.deploymentStatus?.availableReplicas}
            {deployment.deploymentStatus?.updatedReplicas}
          </TableRow>
        ))}
      </Table>
    </DeploymentsContainer>
  );
};

export default DeploymentTable;
