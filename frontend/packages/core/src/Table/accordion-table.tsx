import * as React from "react";
import ChevronRightIcon from "@mui/icons-material/ChevronRight";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import { IconButton as MuiIconButton, TableRow, Theme } from "@mui/material";

import styled from "../styled";

import TableCell from "./components/TableCell";
import type { TableRowProps } from "./types";

const IconButton = styled(MuiIconButton)(({ theme }: { theme: Theme }) => ({
  padding: "0",
  color: theme.palette.secondary[900],
}));

const ChevronRight = styled(ChevronRightIcon)<{ $disabled: boolean }>(
  props => ({ theme }: { theme: Theme }) => ({
    color: props?.$disabled ? theme.palette.secondary[200] : "unset",
  })
);

export interface AccordionRowProps {
  columns: React.ReactElement[];
  children: React.ReactElement<TableRowProps> | React.ReactElement<TableRowProps>[];
}

export const AccordionRow = ({ columns, children }: AccordionRowProps) => {
  const [open, setOpen] = React.useState(false);
  const hasChildren = React.Children.count(children) !== 0;

  const onClick = () => {
    if (hasChildren) {
      setOpen(isOpen => !isOpen);
    }
  };
  return (
    <>
      <TableRow>
        {columns.map((heading: any, index: number) => {
          const icon = (
            <IconButton onClick={onClick} size="large">
              {open ? <KeyboardArrowDownIcon /> : <ChevronRight $disabled={!hasChildren} />}
            </IconButton>
          );
          return (
            <TableCell key={heading} border style={index === 0 ? { display: "flex" } : {}}>
              {index === 0 && icon}
              <div style={{ alignSelf: "center" }}>{heading}</div>
            </TableCell>
          );
        })}
      </TableRow>
      {open && children}
    </>
  );
};
