import * as React from "react";
import styled from "@emotion/styled";
import {
  IconButton,
  Paper as MuiPaper,
  Table as MuiTable,
  TableBody,
  TableCell as MuiTableCell,
  TableContainer as MuiTableContainer,
  TableHead,
  TableProps as MuiTableProps,
  TableRow as MuiTableRow,
  TableRowProps as MuiTableRowProps,
} from "@material-ui/core";
import MoreVertIcon from "@material-ui/icons/MoreVert";

import { Popper, PopperItem } from "../popper";

const TablePaper = styled(MuiPaper)({
  border: "1px solid #E7E7EA",
});

const StyledTableRow = styled(MuiTableRow)({
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

const HeaderActionTableCell = styled(HeaderTableCell)({
  width: "48px",
});

export interface TableContainerProps {
  children: React.ReactElement<MuiTableProps>;
}

export const TableContainer = ({ children }: TableContainerProps) => (
  <MuiTableContainer component={TablePaper} elevation={0}>
    {children}
  </MuiTableContainer>
);

interface TableRowActionProps {
  children: string;
  onClick: () => void;
  icon?: React.ReactElement;
}

export const TableRowAction = ({ children, onClick, icon }: TableRowActionProps) => (
  <PopperItem icon={icon} onClick={onClick}>
    {children}
  </PopperItem>
);

interface TableRowActionsProps {
  children?: React.ReactElement<TableRowActionProps> | React.ReactElement<TableRowActionProps>[];
}

export const TableRowActions = ({ children }: TableRowActionsProps) => {
  const anchorRef = React.useRef(null);
  const [open, setOpen] = React.useState(false);

  return (
    <>
      <IconButton
        disableRipple
        disabled={React.Children.count(children) <= 0}
        ref={anchorRef}
        onClick={() => setOpen(true)}
      >
        <MoreVertIcon />
      </IconButton>
      <Popper open={open} anchorRef={anchorRef} onClickAway={() => setOpen(false)}>
        {children}
      </Popper>
    </>
  );
};

export interface TableProps extends Pick<MuiTableProps, "stickyHeader"> {
  headings?: string[];
  actionsColumn?: boolean;
}

export const Table: React.FC<TableProps> = ({
  headings,
  actionsColumn = false,
  children,
  ...props
}) => {
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
              {actionsColumn && <HeaderActionTableCell />}
            </MuiTableRow>
          </TableHead>
        )}
        <TableBody>{children}</TableBody>
      </MuiTable>
    </TableContainer>
  );
};

export interface TableRowProps extends Pick<MuiTableRowProps, "onClick"> {
  children?: React.ReactNode;
  defaultCellValue?: React.ReactNode;
}

export const TableRow = ({ children = [], onClick, defaultCellValue }: TableRowProps) => (
  <StyledTableRow onClick={onClick}>
    {React.Children.map(children, (value, index) => (
      // eslint-disable-next-line react/no-array-index-key
      <TableCell key={index}>
        {value === null && defaultCellValue !== undefined ? defaultCellValue : value}
      </TableCell>
    ))}
  </StyledTableRow>
);
