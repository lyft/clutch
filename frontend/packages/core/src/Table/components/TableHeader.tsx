import React from "react";
import type { TableSortLabelProps as MuiTableSortLabelProps, Theme } from "@mui/material";
import {
  TableHead as MuiTableHead,
  TableRow as MuiTableRow,
  TableSortLabel as MuiTableSortLabel,
} from "@mui/material";

import styled from "../../styled";
import { Typography } from "../../typography";
import type { Column, TableColumn } from "../types";

import TableCell from "./TableCell";

const StyledTableHeadRow = styled(MuiTableRow)(({ theme }: { theme: Theme }) => ({
  display: "contents",
  backgroundColor: theme.palette.primary[300],
}));

const HeaderCell = styled("div")({
  display: "flex",
  alignItems: "center",
  justifyContent: "space-between",
});

interface TableHeaderProps {
  columns: Column[];
  responsive?: boolean;
  defaultSort?: [string, MuiTableSortLabelProps["direction"]];
  onRequestSort?: (event: React.MouseEvent<unknown>, property: string) => void;
  actionsColumn?: boolean;
  compress?: boolean;
}

const TableHeader = ({
  columns,
  responsive,
  defaultSort: [sortKey, sortDir] = ["", "asc"],
  onRequestSort,
  actionsColumn,
  compress,
}: TableHeaderProps) => {
  const createSortHandler = property => (event: React.MouseEvent<unknown>) => {
    onRequestSort && onRequestSort(event, property);
  };

  const [managedColumns, setManagedColumns] = React.useState<TableColumn[]>([]);

  React.useEffect(() => {
    if (columns?.length === 0) {
      // eslint-disable-next-line no-console
      console.warn("Table must have at least one column.");
    } else {
      setManagedColumns(
        columns.map((c, i) => {
          if (React.isValidElement(c)) {
            return { id: `element${i}`, render: c };
          }
          if (typeof c === "string") {
            return { id: c, title: c };
          }
          return c as TableColumn;
        })
      );
    }
  }, [columns]);

  return (
    managedColumns?.length !== 0 &&
    managedColumns.filter(h => h?.title?.length !== 0 || h?.render).length !== 0 && (
      <MuiTableHead>
        <StyledTableHeadRow>
          {managedColumns.map(h => (
            <TableCell
              key={h?.id}
              responsive={responsive}
              align="left"
              sortDirection={h?.sortable && sortKey === h?.id ? sortDir : false}
            >
              {h?.sortable ? (
                <HeaderCell>
                  <MuiTableSortLabel
                    active={sortKey === h?.id}
                    direction={sortKey === h?.id ? sortDir : "asc"}
                    onClick={createSortHandler(h?.id)}
                  >
                    <Typography variant="subtitle3">{h?.title}</Typography>
                  </MuiTableSortLabel>
                  {h?.options}
                </HeaderCell>
              ) : (
                <HeaderCell>
                  <Typography variant="subtitle3">{h?.title || h?.render}</Typography>
                  {h?.options}
                </HeaderCell>
              )}
            </TableCell>
          ))}
          {actionsColumn && !(responsive && compress) && (
            <TableCell responsive={responsive} action />
          )}
        </StyledTableHeadRow>
      </MuiTableHead>
    )
  );
};

export default TableHeader;
