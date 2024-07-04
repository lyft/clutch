import type {
  TableCellProps as MuiTableCellProps,
  TableProps as MuiTableProps,
  TableRowProps as MuiTableRowProps,
  TableSortLabelProps as MuiTableSortLabelProps,
} from "@mui/material";
import type { Breakpoint } from "@mui/material/styles";

interface TableRowProps
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

interface TableColumn {
  id: string;
  title?: string;
  sortable?: boolean;
  render?: JSX.Element;
  filter?: boolean;
  filterRender?: JSX.Element;
}

type Column = string | TableColumn | JSX.Element;
interface TableProps extends Pick<MuiTableProps, "stickyHeader"> {
  /** The names of the columns. This must be set (even to empty string) to render the table. */
  columns: Column[];
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
  defaultSort?: [string, MuiTableSortLabelProps["direction"]];
  onRequestSort?: (event: React.MouseEvent<unknown>, property: string) => void;
}

interface TableContainerProps {
  children: React.ReactElement<TableProps>;
}

interface TableCellProps extends MuiTableCellProps, Pick<TableProps, "responsive"> {
  action?: boolean;
  border?: boolean;
}

interface TableHeaderProps {
  columns: Column[];
  responsive?: boolean;
  defaultSort?: [string, MuiTableSortLabelProps["direction"]];
  onRequestSort?: (event: React.MouseEvent<unknown>, property: string) => void;
  actionsColumn?: boolean;
  compress?: boolean;
}

export type {
  TableRowProps,
  TableColumn,
  TableProps,
  TableContainerProps,
  TableCellProps,
  TableHeaderProps,
};
