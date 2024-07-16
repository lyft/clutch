import * as React from "react";
import EmojiPeopleIcon from "@mui/icons-material/EmojiPeople";
import FilterListIcon from "@mui/icons-material/FilterList";
import { action } from "@storybook/addon-actions";
import type { Meta } from "@storybook/react";

import { Table, TableRow, TableRowAction, TableRowActions } from "../table";
import type { TableProps, TableRowProps } from "../types";

export default {
  title: "Core/Table/Table",
  component: Table,
} as Meta;

const Template = ({ row, ...props }: TableProps & { row: React.ReactElement }) => (
  <div style={{ maxHeight: "300px", display: "flex" }}>
    <Table {...props} columns={["Column 1", "Column 2", "Column 3", "Column 4", "Column 5"]}>
      {Array(10)
        .fill(null)
        // eslint-disable-next-line react/no-array-index-key
        .map((_, index: number) => React.cloneElement(row, { key: index }))}
    </Table>
  </div>
);

const Template2 = ({ row, ...props }: TableProps & { row: React.ReactElement }) => (
  <div style={{ maxHeight: "300px", display: "flex" }}>
    <Table {...props}>
      {Array(10)
        .fill(null)
        // eslint-disable-next-line react/no-array-index-key
        .map((_, index: number) => React.cloneElement(row, { key: index }))}
    </Table>
  </div>
);

const PrimaryTableRow = (props: TableRowProps) => (
  <TableRow {...props}>
    <div>Value 1</div>
    <div>Value 2</div>
    <div>Value 3</div>
    <div>Value 4</div>
    <div>Value 5</div>
  </TableRow>
);

const IncompleteTableRow = (props: TableRowProps) => {
  let data;
  return (
    <TableRow {...props}>
      <div>Value 1</div>
      {data}
      {data}
      {data}
      {data}
    </TableRow>
  );
};

const ActionableTableRow = (props: TableRowProps) => (
  <TableRow {...props}>
    <div>Value 1</div>
    <div>Value 2</div>
    <div>Value 3</div>
    <div>Value 4</div>
    <div>Value 5</div>
    <TableRowActions>
      <TableRowAction icon={<EmojiPeopleIcon />} onClick={action("row-action")}>
        Take Action
      </TableRowAction>
    </TableRowActions>
  </TableRow>
);

export const Primary = Template.bind({});
Primary.args = {
  row: <PrimaryTableRow />,
};

export const StickyHeader = Template.bind({});
StickyHeader.args = {
  stickyHeader: true,
  row: <PrimaryTableRow />,
};

export const MissingCellValue = Template.bind({});
MissingCellValue.args = {
  row: <IncompleteTableRow />,
};

export const DefaultCellValue = Template.bind({});
DefaultCellValue.args = {
  row: <IncompleteTableRow cellDefault="null" />,
};

export const ActionableRows = Template.bind({});
ActionableRows.args = {
  actionsColumn: true,
  row: <ActionableTableRow />,
};

export const SortFilterOptions = Template2.bind({});
SortFilterOptions.args = {
  defaultSort: ["column1", "asc"],
  onRequestSort: () => {},
  columns: [
    {
      id: "column1",
      title: "Column 1",
      sortable: true,
      options: <FilterListIcon fontSize="medium" />,
    },
    {
      id: "column2",
      title: "Column 2",
      sortable: true,
      options: <FilterListIcon fontSize="medium" />,
    },
    {
      id: "column3",
      title: "Column 3",
      sortable: true,
      options: <FilterListIcon fontSize="medium" />,
    },
    {
      id: "column4",
      title: "Column 4",
      sortable: true,
      options: <FilterListIcon fontSize="medium" />,
    },
    {
      id: "column5",
      title: "Column 5",
      sortable: true,
      options: <FilterListIcon fontSize="medium" />,
    },
  ],
  row: <PrimaryTableRow />,
};
