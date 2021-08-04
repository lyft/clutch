import * as React from "react";
import styled from "@emotion/styled";
import {
  IconButton,
  Paper as MuiPaper,
  Table as MuiTable,
  TableBody as MuiTableBody,
  TableCell as MuiTableCell,
  TableContainer as MuiTableContainer,
  TableHead as MuiTableHead,
  TableProps as MuiTableProps,
  TableRow as MuiTableRow,
  TableRowProps as MuiTableRowProps,
  useMediaQuery,
} from "@material-ui/core";
import MoreVertIcon from "@material-ui/icons/MoreVert";

import { Popper, PopperItem } from "../popper";
import { Typography } from "../typography";

const StyledPaper = styled(MuiPaper)({
  border: "1px solid #E7E7EA",
});

const StyledTable = styled(MuiTable)<{ actions: boolean; columnCount: number; compress: boolean }>(
  {
    minWidth: "100%",
    borderCollapse: "collapse",
    alignItems: "center",
  },
  props => ({
    display: props.compress ? "block" : "grid",
    gridTemplateColumns: `repeat(${props.columnCount}, auto)${props.actions ? " 80px" : ""}`,
  })
);

const StyledTableBody = styled(MuiTableBody)({
  display: "contents",
});

const StyledTableHead = styled(MuiTableHead)({
  display: "contents",
  backgroundColor: "#D7DAF6",
});

const StyledTableRow = styled(MuiTableRow)({
  display: "contents",
  ":nth-child(even)": {
    background: "rgba(13, 16, 48, 0.03)",
  },
  ":hover": {
    background: "#EBEDFB",
  },
});

const StyledTableCell = styled(MuiTableCell)<{ border?: boolean }>(
  {
    display: "flex",
    alignItems: "center",
    fontSize: "14px",
    padding: "15px 16px",
    color: "#0D1030",
    overflow: "hidden",
    background: "inherit",
    minHeight: "100%",
  },
  props => ({
    borderBottom: props?.border ? "1px solid #E7E7EA" : "0",
  })
);

export interface TableContainerProps {
  children: React.ReactElement<MuiTableProps>;
}

const TableContainer = ({ children }: TableContainerProps) => (
  <MuiTableContainer component={StyledPaper} elevation={0}>
    {children}
  </MuiTableContainer>
);

export interface TableProps extends Pick<MuiTableProps, "stickyHeader"> {
  /** The names of the columns. This must be set (even to empty string) to render the table. */
  columns: string[];
  /** Hide the header. By default this is false. */
  hideHeader?: boolean;
  /** Add an actions column. By default this is false. */
  actionsColumn?: boolean;
}

const Table: React.FC<TableProps> = ({
  columns,
  hideHeader = false,
  actionsColumn = false,
  children,
  ...props
}) => {
  const showHeader = !hideHeader;
  const compress = useMediaQuery((theme: any) => theme.breakpoints.down("md"));

  return (
    <TableContainer>
      <StyledTable
        compress={compress}
        columnCount={columns?.length}
        actions={actionsColumn}
        {...props}
      >
        {/*
          Filter out empty strings from column headers.
          This may be unintended which is why we override wit hthe hideHeader prop.
        */}
        {(showHeader ||
          (columns?.length !== 0 && columns.filter(h => h.length !== 0).length !== 0)) && (
          <StyledTableHead>
            {columns.map(h => (
              <StyledTableCell>
                <Typography variant="subtitle3">{h}</Typography>
              </StyledTableCell>
            ))}
            {actionsColumn && !compress && <StyledTableCell />}
          </StyledTableHead>
        )}
        <StyledTableBody>{children}</StyledTableBody>
      </StyledTable>
    </TableContainer>
  );
};

export interface TableRowProps extends Pick<MuiTableRowProps, "onClick"> {
  children?: React.ReactNode;
  /**
   * The default element to render if children are null. If not present and a child is null
   * this the child's value will be used.
   */
  cellDefault?: React.ReactNode;
}

const TableRow = ({ children = [], onClick, cellDefault }: TableRowProps) => (
  <StyledTableRow onClick={onClick}>
    {React.Children.map(children, (value, index) => (
      // eslint-disable-next-line react/no-array-index-key
      <StyledTableCell key={index}>
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

export {
  StyledTableCell as TableCell,
  Table,
  TableContainer,
  TableRow,
  TableRowAction,
  TableRowActions,
};
