import React from "react";
import type { Theme } from "@mui/material";
import { TableCell as MuiTableCell } from "@mui/material";

import styled from "../../styled";
import type { TableCellProps } from "../types";

const StyledTableCell = styled(MuiTableCell)<{
  $border?: TableCellProps["border"];
  $responsive?: TableCellProps["responsive"];
  $action?: TableCellProps["action"];
}>(
  ({ theme }: { theme: Theme }) => ({
    alignItems: "center",
    fontSize: "14px",
    padding: "15px 16px",
    color: theme.palette.secondary[900],
    overflow: "hidden",
    background: "inherit",
    minHeight: "100%",
  }),
  props => ({ theme }: { theme: Theme }) => ({
    borderBottom: props?.$border ? `1px solid ${theme.palette.secondary[200]}` : "0",
    display: props.$responsive ? "flex" : "",
    width: !props.$responsive && props.$action ? "80px" : "",
  })
);

const TableCell = ({ action, border, responsive, ...props }: TableCellProps) => (
  <StyledTableCell $action={action} $border={border} $responsive={responsive} {...props} />
);

export default TableCell;
