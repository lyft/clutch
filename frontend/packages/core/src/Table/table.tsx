import * as React from "react";
import MoreVertIcon from "@mui/icons-material/MoreVert";
import type { Theme } from "@mui/material";
import {
  IconButton,
  Paper as MuiPaper,
  Table as MuiTable,
  TableBody as MuiTableBody,
  TableContainer as MuiTableContainer,
  TableRow as MuiTableRow,
  useMediaQuery,
} from "@mui/material";

import { Popper, PopperItem } from "../popper";
import styled from "../styled";

import TableCell from "./components/TableCell";
import TableHeader from "./components/TableHeader";
import type { TableContainerProps, TableProps, TableRowProps } from "./types";

const StyledPaper = styled(MuiPaper)(({ theme }: { theme: Theme }) => ({
  border: `1px solid ${theme.palette.secondary[200]}`,
}));

const StyledTable = styled(MuiTable)<{
  $hasActionsColumn?: TableProps["actionsColumn"];
  $columnCount: number;
  $compress: boolean;
  $responsive?: TableProps["responsive"];
  $overflow?: TableProps["overflow"];
}>(
  {
    minWidth: "100%",
    borderCollapse: "collapse",
    alignItems: "center",
  },
  props => ({
    display: !props.$responsive ? "table" : props.$compress ? "table" : "grid",
    gridTemplateColumns: `repeat(${props.$columnCount}, auto)${
      props.$hasActionsColumn ? " 80px" : ""
    }`,
    ".MuiTableCell-root": {
      wordBreak: props.$overflow === "scroll" ? "normal" : props.$overflow,
    },
  })
);

const StyledTableBody = styled(MuiTableBody)({
  display: "contents",
});

const StyledTableRow = styled(MuiTableRow)<{
  $responsive?: TableRowProps["responsive"];
}>(
  ({ theme }: { theme: Theme }) => ({
    ":nth-of-type(even)": {
      background: theme.palette.secondary[50],
    },
    ":hover": {
      background: theme.palette.primary[200],
    },
  }),
  props => ({
    display: props.$responsive ? "contents" : "",
  })
);

const TableContainer = ({ children }: TableContainerProps) => (
  <MuiTableContainer component={StyledPaper} elevation={0}>
    {children}
  </MuiTableContainer>
);

const Table: React.FC<TableProps> = React.forwardRef(
  (
    {
      columns,
      compressBreakpoint = "sm",
      hideHeader = false,
      actionsColumn = false,
      responsive = false,
      overflow = "scroll",
      children,
      defaultSort,
      onRequestSort,
      ...props
    },
    ref
  ) => {
    const showHeader = !hideHeader;
    const compress = useMediaQuery((theme: any) => theme.breakpoints.down(compressBreakpoint));

    return (
      <TableContainer>
        <StyledTable
          $compress={compress}
          $columnCount={columns?.length}
          $hasActionsColumn={actionsColumn}
          $responsive={responsive}
          $overflow={overflow}
          ref={ref}
          {...props}
        >
          {/*
          Filter out empty strings from column headers.
          This may be unintended which is why we override with the hideHeader prop.
        */}
          {showHeader && (
            <TableHeader
              columns={columns}
              responsive={responsive}
              defaultSort={defaultSort}
              onRequestSort={onRequestSort}
              actionsColumn={actionsColumn}
              compress={compress}
            />
          )}
          <StyledTableBody>
            {React.Children.map(children, (c: React.ReactElement<TableRowProps>) =>
              React.cloneElement(c, { responsive })
            )}
          </StyledTableBody>
        </StyledTable>
      </TableContainer>
    );
  }
);

const TableRow = ({
  children = [],
  onClick,
  cellDefault,
  responsive = false,
  colSpan,
  ...props
}: TableRowProps) => (
  <StyledTableRow onClick={onClick} $responsive={responsive} {...props}>
    {React.Children.map(children, (value, index) => (
      // eslint-disable-next-line react/no-array-index-key
      <TableCell key={index} responsive={responsive} colSpan={colSpan}>
        {value === null && cellDefault !== undefined ? cellDefault : value}
      </TableCell>
    ))}
  </StyledTableRow>
);

interface TableRowActionProps {
  children: string;
  onClick: () => void;
  icon?: React.ReactElement;
}

const TableRowAction = ({ children, onClick, icon }: TableRowActionProps) => (
  <PopperItem icon={icon} onClick={onClick}>
    {children}
  </PopperItem>
);

interface TableRowActionsProps {
  children?: React.ReactElement<TableRowActionProps> | React.ReactElement<TableRowActionProps>[];
}

const TableRowActions = ({ children }: TableRowActionsProps) => {
  const anchorRef = React.useRef(null);
  const [open, setOpen] = React.useState(false);

  return (
    <>
      <IconButton
        disableRipple
        disabled={React.Children.count(children) <= 0}
        ref={anchorRef}
        onClick={() => setOpen(true)}
        size="large"
      >
        <MoreVertIcon />
      </IconButton>
      <Popper
        open={open}
        anchorRef={anchorRef}
        onClickAway={() => setOpen(false)}
        placement="bottom-end"
      >
        {children}
      </Popper>
    </>
  );
};

export { TableCell, Table, TableContainer, TableRow, TableRowAction, TableRowActions };
