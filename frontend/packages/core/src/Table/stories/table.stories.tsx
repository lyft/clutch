import * as React from "react";
import EmojiPeopleIcon from "@material-ui/icons/EmojiPeople";
import { action } from "@storybook/addon-actions";
import type { Meta } from "@storybook/react";

import {
  Table,
  TableProps,
  TableRow,
  TableRowAction,
  TableRowActions,
  TableRowProps,
} from "../table";

export default {
  title: "Core/Table/Table",
  component: Table,
} as Meta;

const Template = ({ row, ...props }: TableProps & { row: React.ReactElement }) => (
  <div style={{ maxHeight: "300px", display: "flex" }}>
    <Table headings={["Column 1", "Column 2"]} {...props}>
      {
        // eslint-disable-next-line react/no-array-index-key
        [...Array(10)].map((_, index: number) => React.cloneElement(row, { key: index }))
      }
    </Table>
  </div>
);

const PrimaryTableRow = (props: TableRowProps) => (
  <TableRow {...props}>
    <div>Value 1</div>
    <div>Value 2</div>
  </TableRow>
);

const IncompleteTableRow = (props: TableRowProps) => {
  let data;
  return (
    <TableRow {...props}>
      <div>Value 1</div>
      {data}
    </TableRow>
  );
};

const ActionableTableRow = (props: TableRowProps) => (
  <TableRow {...props}>
    <div>Value 1</div>
    <div>Value 2</div>
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
  row: <IncompleteTableRow defaultCellValue="null" />,
};

export const ActionableRows = Template.bind({});
ActionableRows.args = {
  actionsColumn: true,
  row: <ActionableTableRow />,
};
