import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import Paper from "@material-ui/core/Paper";
import { makeStyles } from "@material-ui/core/styles";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableContainer from "@material-ui/core/TableContainer";
import TableHead from "@material-ui/core/TableHead";
import TablePagination from "@material-ui/core/TablePagination";
import TableRow from "@material-ui/core/TableRow";
import TableSortLabel from "@material-ui/core/TableSortLabel";

import { compareProperties, propertyToString } from "./property-helpers";

const useStyles = makeStyles(theme => ({
  root: {
    width: "100%",
  },
  paper: {
    width: "100%",
    marginBottom: theme.spacing(2),
  },
  table: {
    minWidth: 750,
  },
}));

type Ordering = "asc" | "desc";
type ListViewItem = IClutch.chaos.experimentation.v1.ListViewItem;

function getComparator(order: Ordering, orderBy: string) {
  return order === "desc"
    ? (a: ListViewItem, b: ListViewItem) =>
        compareProperties(a.properties.items[orderBy], b.properties.items[orderBy])
    : (a: ListViewItem, b: ListViewItem) =>
        -compareProperties(a.properties.items[orderBy], b.properties.items[orderBy]);
}

function stableSort(
  array: ListViewItem[],
  comparator: (a: ListViewItem, b: ListViewItem) => number
): ListViewItem[] {
  const stabilizedThis = array.map((el, index) => {
    return { item: el, index };
  });
  stabilizedThis.sort((a, b) => {
    const order = comparator(a.item, b.item);
    if (order !== 0) {
      return order;
    }
    return a.index - b.index;
  });
  return stabilizedThis.map(el => el.item);
}

interface EnhancedTableHeadProps {
  columns: Column[];
  order: Ordering;
  orderBy: string;
  onRequestSort: (event: any, columnIdentifier: string) => void;
}

const EnhancedTableHead: React.FC<EnhancedTableHeadProps> = ({
  columns,
  order,
  orderBy,
  onRequestSort,
}) => {
  const createSortHandler = property => event => {
    onRequestSort(event, property);
  };

  return (
    <TableHead>
      <TableRow>
        {columns.map(column => (
          <TableCell
            key={column.identifier}
            align="left"
            sortDirection={orderBy === column.identifier ? order : false}
          >
            {column.sortable && (
              <TableSortLabel
                active={orderBy === column.identifier}
                direction={orderBy === column.identifier ? order : "asc"}
                onClick={createSortHandler(column.identifier)}
              >
                {column.header}
              </TableSortLabel>
            )}
            {!column.sortable && column.header}
          </TableCell>
        ))}
      </TableRow>
    </TableHead>
  );
};

interface Column {
  identifier: string;
  header: string;
  sortable: boolean;
}

interface ListViewProps {
  columns: Column[];
  items: ListViewItem[];
  onRowSelection: (event: any, item: ListViewItem) => void;
}

const ListView: React.FC<ListViewProps> = ({ columns, items, onRowSelection }) => {
  const classes = useStyles();
  const [order, setOrder] = React.useState<Ordering | undefined>("asc");
  const [orderBy, setOrderBy] = React.useState("");
  const [page, setPage] = React.useState(0);
  const [rowsPerPage, setRowsPerPage] = React.useState(25);

  const handleRequestSort = (event, columnIdentifier) => {
    const isAsc = orderBy === columnIdentifier && order === "asc";
    setOrder(isAsc ? "desc" : "asc");
    setOrderBy(columnIdentifier);
  };

  const handleClick = (event: any, item: ListViewItem) => {
    onRowSelection(event, item);
  };

  const handleChangePage = (event: any, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = event => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  return (
    <div className={classes.root}>
      <Paper className={classes.paper}>
        <TableContainer>
          <Table className={classes.table} size="medium">
            <EnhancedTableHead
              columns={columns}
              order={order}
              orderBy={orderBy}
              onRequestSort={handleRequestSort}
            />
            {items && (
              <TableBody>
                {stableSort(items, getComparator(order, orderBy))
                  .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                  .map((item: ListViewItem) => {
                    return (
                      <TableRow
                        hover
                        onClick={event => handleClick(event, item)}
                        key={item.identifier.toString()}
                      >
                        {columns &&
                          columns.map(column => {
                            return (
                              <TableCell key={column.identifier} align="left">
                                {propertyToString(item.properties.items[column.identifier])}
                              </TableCell>
                            );
                          })}
                      </TableRow>
                    );
                  })}
              </TableBody>
            )}
          </Table>
        </TableContainer>
        <TablePagination
          rowsPerPageOptions={[25, 50, 100]}
          component="div"
          count={items?.length ?? 0}
          rowsPerPage={rowsPerPage}
          page={page}
          onChangePage={handleChangePage}
          onChangeRowsPerPage={handleChangeRowsPerPage}
        />
      </Paper>
    </div>
  );
};

export { Column, ListView };
