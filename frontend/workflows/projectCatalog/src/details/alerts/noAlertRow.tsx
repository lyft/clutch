import React from "react";
import { Chip, TableRow, Typography } from "@clutch-sh/core";

const NoAlertRow = ({ project }: { project: string }) => {
  return (
    <TableRow key={project}>
      <Typography variant="subtitle3">{project}</Typography>
      <Chip label="No Alerts" variant="neutral" />
      <></>
      <></>
    </TableRow>
  );
};

export default NoAlertRow;
