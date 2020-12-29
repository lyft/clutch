import React from "react";
import { Collapse, IconButton, TableCell, TableRow } from "@material-ui/core";
import KeyboardArrowDownIcon from "@material-ui/icons/KeyboardArrowDown";
import KeyboardArrowUpIcon from "@material-ui/icons/KeyboardArrowUp";

import { Status } from "../icon";

import { Table } from "./table";

interface StatusRowProps {
  success: boolean;
  data: any[];
}

const StatusRow: React.FC<StatusRowProps> = ({ success, data }) => {
  const displayData = [...data];
  const headerValue = displayData.shift();
  const variant = success ? "success" : "failure";
  return (
    <TableRow>
      <TableCell align="left">
        <Status variant={variant}>{headerValue}</Status>
      </TableCell>
      {displayData.map(value => (
        <TableCell key={value} align="left">
          {value}
        </TableCell>
      ))}
    </TableRow>
  );
};

interface ExpandableTableProps {
  headings: string[];
}

const ExpandableTable: React.FC<ExpandableTableProps> = ({ headings, children }) => (
  <Table stickyHeader headings={[...headings, ""]}>
    {children}
  </Table>
);

interface ExpandableRowProps {
  heading: string;
  summary: React.ReactElement | string;
}

const ExpandableRow: React.FC<ExpandableRowProps> = ({ heading, summary, children }) => {
  const [open, setOpen] = React.useState(false);
  return (
    <>
      <TableRow>
        <TableCell component="th" scope="row">
          <strong>{heading}</strong>
        </TableCell>
        {summary !== undefined && <TableCell align="right">{summary}</TableCell>}
        <TableCell align="center">
          {React.Children.toArray(children).length !== 0 && (
            <IconButton aria-label="expand row" size="small" onClick={() => setOpen(!open)}>
              {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
            </IconButton>
          )}
        </TableCell>
      </TableRow>
      <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <Table>{children}</Table>
          </Collapse>
        </TableCell>
      </TableRow>
    </>
  );
};

export { ExpandableRow, ExpandableTable, StatusRow };
