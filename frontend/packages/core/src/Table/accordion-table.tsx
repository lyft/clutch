import * as React from "react";
import ChevronRightIcon from "@mui/icons-material/ChevronRight";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import { IconButton as MuiIconButton, TableRow } from "@mui/material";

import styled from "../styled";

import type { TableRowProps } from "./table";
import { TableCell } from "./table";

const IconButton = styled(MuiIconButton)({
  padding: "0",
  color: "#0D1030",
});

const ChevronRight = styled(ChevronRightIcon)<{ $disabled: boolean }>(props => ({
  color: props?.$disabled ? "#E7E7EA" : "unset",
}));

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
