import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Table, TableRow } from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import styled from "@emotion/styled";
import _ from "lodash";

const CronsContainer = styled.div({
  display: "flex",
  maxHeight: "50vh",
});

const CronTable = () => {
  const cronListData = useDataLayout("cronListData", { hydrate: false });
  const crons = cronListData.displayValue()?.cronJobs as IClutch.k8s.v1.CronJob[];

  return (
    <CronsContainer>
      <Table
        stickyHeader
        actionsColumn
        headings={["Name", "Cluster", "Schedule", "Suspend", "Active Jobs", "Concurrency Policy"]}
      >
        {_.sortBy(crons, [
          o => {
            return o.name;
          },
        ]).map(cron => (
          <TableRow key={cron.name} defaultCellValue="nil">
            {cron.name}
            {cron.cluster}
            {cron.schedule}
            {cron.suspend}
            {cron.numActiveJobs}
            {cron.concurrencyPolicy}
          </TableRow>
        ))}
      </Table>
    </CronsContainer>
  );
};

export default CronTable;
