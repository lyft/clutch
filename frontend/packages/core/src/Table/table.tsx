import * as React from "react";
import styled from "@emotion/styled";
import type { TableProps as MuiTableProps } from "@material-ui/core";
import {
  Paper,
  Table as MuiTable,
  TableBody,
  TableCell as MuiTableCell,
  TableContainer as MuiTableContainer,
  TableHead,
  TableRow as MuiTableRow,
} from "@material-ui/core";

const TablePaper = styled(Paper)({
  border: "1px solid #E7E7EA",
});

export const StyledTableRow = styled(MuiTableRow)({
  ":hover": {
    background: "#EBEDFB",
  },
});

export const TableCell = styled(MuiTableCell)(
  {
    fontSize: "14px",
    padding: "12px 16px",
    color: "#0D1030",
  },
  props => ({
    borderBottom: props["data-border"] ? "1px solid #E7E7EA" : "0",
  })
);

const HeaderTableCell = styled(TableCell)({
  backgroundColor: "rgba(248, 248, 249, 1)",
  fontWeight: 600,
});

export interface TableContainerProps {
  children: React.ReactElement<MuiTableProps>;
}

export const TableContainer = ({ children }: TableContainerProps) => (
  <MuiTableContainer component={TablePaper} elevation={0}>
    {children}
  </MuiTableContainer>
);

export interface TableProps extends Pick<MuiTableProps, "stickyHeader"> {
  headings?: string[];
}

export const Table: React.FC<TableProps> = ({ headings, children, ...props }) => {
  const localHeadings = headings ? [...headings] : [];

  return (
    <TableContainer>
      <MuiTable {...props}>
        {localHeadings.length !== 0 && (
          <TableHead>
            <MuiTableRow>
              {localHeadings.map(heading => (
                <HeaderTableCell key={heading} align="left">
                  {heading}
                </HeaderTableCell>
              ))}
            </MuiTableRow>
          </TableHead>
        )}
        <TableBody>{children}</TableBody>
      </MuiTable>
    </TableContainer>
  );
};

export interface TableRowProps {
  children?: React.ReactNode;
}

export const TableRow = ({ children = [] }: TableRowProps) => (
  <StyledTableRow>
    {React.Children.map(children, (value, index) => (
      // eslint-disable-next-line react/no-array-index-key
      <TableCell key={index}>{value}</TableCell>
    ))}
  </StyledTableRow>
);
