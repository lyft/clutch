import * as React from "react";
import MoreVertIcon from "@mui/icons-material/MoreVert";
import type {
  TableCellProps as MuiTableCellProps,
  TableProps as MuiTableProps,
  TableRowProps as MuiTableRowProps,
  Theme,
} from "@mui/material";
import {
  IconButton,
  Paper as MuiPaper,
  Table as MuiTable,
  TableBody as MuiTableBody,
  TableCell as MuiTableCell,
  TableContainer as MuiTableContainer,
  TableHead as MuiTableHead,
  TableRow as MuiTableRow,
  useMediaQuery,
} from "@mui/material";
import type { Breakpoint } from "@mui/material/styles";

import { Popper, PopperItem } from "../popper";
import styled from "../styled";
import { Typography } from "../typography";

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

const StyledTableHeadRow = styled(MuiTableRow)(({ theme }: { theme: Theme }) => ({
  display: "contents",
  backgroundColor: theme.palette.primary[300],
}));

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

interface TableCellProps extends MuiTableCellProps, Pick<TableProps, "responsive"> {
  action?: boolean;
  border?: boolean;
}

const TableCell = ({ action, border, responsive, ...props }: TableCellProps) => (
  <StyledTableCell $action={action} $border={border} $responsive={responsive} {...props} />
);

interface TableContainerProps {
  children: React.ReactElement<TableProps>;
}

const TableContainer = ({ children }: TableContainerProps) => (
  <MuiTableContainer component={StyledPaper} elevation={0}>
    {children}
  </MuiTableContainer>
);

interface TableProps extends Pick<MuiTableProps, "stickyHeader"> {
  /** The names of the columns. This must be set (even to empty string) to render the table. */
  columns: string[];
  /** The breakpoint at which to compress the table rows. By default the small breakpoint is used. */
  compressBreakpoint?: Breakpoint;
  /** Hide the header. By default this is false. */
  hideHeader?: boolean;
  /** Add an actions column. By default this is false. */
  actionsColumn?: boolean;
  /** Make table responsive */
  responsive?: boolean;
  /** How to handle horizontal overflow */
  overflow?: "scroll" | "break-word";
  /** Table rows to render */
  children?:
    | (React.ReactElement<TableRowProps> | null | undefined | {})[]
    | React.ReactElement<TableRowProps>;
}

const Table: React.FC<TableProps> = ({
  columns,
  compressBreakpoint = "sm",
  hideHeader = false,
  actionsColumn = false,
  responsive = false,
  overflow = "scroll",
  children,
  ...props
}) => {
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
        {...props}
      >
        {/*
          Filter out empty strings from column headers.
          This may be unintended which is why we override wit hthe hideHeader prop.
        */}
        {showHeader && columns?.length !== 0 && columns.filter(h => h.length !== 0).length !== 0 && (
          <MuiTableHead>
            <StyledTableHeadRow>
              {columns.map(h => (
                <StyledTableCell key={h} $responsive={responsive}>
                  <Typography variant="subtitle3">{h}</Typography>
                </StyledTableCell>
              ))}
              {actionsColumn && !(responsive && compress) && (
                <StyledTableCell $responsive={responsive} $action />
              )}
            </StyledTableHeadRow>
          </MuiTableHead>
        )}
        <StyledTableBody>
          {React.Children.map(children, (c: React.ReactElement<TableRowProps>) =>
            React.cloneElement(c, { responsive })
          )}
        </StyledTableBody>
      </StyledTable>
    </TableContainer>
  );
};

export interface TableRowProps
  extends Pick<MuiTableRowProps, "onClick">,
    Pick<MuiTableCellProps, "colSpan"> {
  children?: React.ReactNode;
  /**
   * The default element to render if children are null. If not present and a child is null
   * this the child's value will be used.
   */
  cellDefault?: React.ReactNode;
  /**
   * Make the table row responsive. This is mainly used for internal rendering. Consumers
   * should set the responsive prop on the table.
   */
  responsive?: boolean;
}

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
      <StyledTableCell key={index} $responsive={responsive} colSpan={colSpan}>
        {value === null && cellDefault !== undefined ? cellDefault : value}
      </StyledTableCell>
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

export type { TableProps };
export type { TableContainerProps };
