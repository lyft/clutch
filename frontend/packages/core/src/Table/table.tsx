import React from "react";
import type { 
  TableProps as MuiTableProps,
  TableRowProps
} from "@material-ui/core";
import {
  Paper,
  Table as MuiTable,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from "@material-ui/core";
import styled from "styled-components";

const TablePaper = styled(Paper)`
  max-height: 400px;
`;

interface TableProps extends MuiTableProps {
  headings?: string[];
  elevation?: number;
}

const Table: React.FC<TableProps> = ({ headings, children, elevation = 1, ...props }) => {
  let localHeadings = [];
  if (headings) {
    localHeadings = [...headings];
  }

  return (
    // n.b. material ui doesn't use the component prop to determine the prop types.
    // @ts-ignore
    <TableContainer component={TablePaper} elevation={elevation}>
      <MuiTable {...props}>
        {localHeadings.length !== 0 && (
          <TableHead>
            <TableRow>
              <TableCell align="left">
                <strong>{localHeadings.shift()}</strong>
              </TableCell>
              {localHeadings.map(heading => (
                <TableCell key={heading} align="right">
                  <strong>{heading}</strong>
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
        )}
        <TableBody>{children}</TableBody>
      </MuiTable>
    </TableContainer>
  );
};

interface RowProps extends TableRowProps {
  data: any[];
}

const Row: React.FC<RowProps> = ({ data, ...props }) => {
  const rowData = data ? [...data] : [];
  const headerValue = rowData.shift();
  return (
    <TableRow {...props}>
      <TableCell component="th" scope="row">
        {headerValue}
      </TableCell>
      {rowData.map(value => (
        <TableCell key={value} align="right">
          {value}
        </TableCell>
      ))}
    </TableRow>
  );
};

export { Row, Table, TablePaper };
