import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Table, TableRow } from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import styled from "@emotion/styled";
import _ from "lodash";

const StatefulSetsContainer = styled.div({
  display: "flex",
  maxHeight: "50vh",
});

const StatefulSetTable = () => {
  const statefulSetListData = useDataLayout("statefulSetListData", { hydrate: false });
  const statefulSets = statefulSetListData.displayValue()
    ?.statefulSets as IClutch.k8s.v1.StatefulSet[];

  return (
    <StatefulSetsContainer>
      <Table
        stickyHeader
        actionsColumn
        columns={["Name", "Cluster", "Replicas Desired", "Replicas Ready", "Replicas Up-To-Date"]}
      >
        {_.sortBy(statefulSets, [
          o => {
            return o.name;
          },
        ]).map(statefulSet => (
          <TableRow key={statefulSet.name} cellDefault="nil">
            {statefulSet.name}
            {statefulSet.cluster}
            {statefulSet.status?.replicas}
            {statefulSet.status?.readyReplicas}
            {statefulSet.status?.updatedReplicas}
          </TableRow>
        ))}
      </Table>
    </StatefulSetsContainer>
  );
};

export default StatefulSetTable;
