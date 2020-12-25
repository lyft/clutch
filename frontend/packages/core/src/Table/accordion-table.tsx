import * as React from "react";
import styled from "@emotion/styled";
import { IconButton as MuiIconButton, TableRow } from "@material-ui/core";
import ChevronRightIcon from "@material-ui/icons/ChevronRight";
import KeyboardArrowDownIcon from "@material-ui/icons/KeyboardArrowDown";

import { TableCell } from "./table";

const IconButton = styled(MuiIconButton)({
  padding: "0",
});

const ChevronRight = styled(ChevronRightIcon)(props => ({
  color: props?.["data-disabled"] ? "#E7E7EA" : "unset",
}));

export interface AccordionRowProps {
  headings: React.ReactElement[];
}

export const AccordionRow: React.FC<AccordionRowProps> = ({ headings, children }) => {
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
        {headings.map((heading: any, index: number) => {
          const icon = (
            <IconButton onClick={onClick}>
              {open ? <KeyboardArrowDownIcon /> : <ChevronRight data-disabled={!hasChildren} />}
            </IconButton>
          );
          return (
            <TableCell key={heading} style={index === 0 ? { display: "flex" } : {}}>
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
