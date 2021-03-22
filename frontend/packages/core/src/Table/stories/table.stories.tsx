import * as React from "react";
import type { Meta } from "@storybook/react";

import type { TableProps, TableRowProps } from "../table";
import { Table, TableRow } from "../table";

export default {
  title: "Core/Table/Table",
  component: Table,
} as Meta;

const Template = ({ row, ...props }: TableProps & { row: React.ReactElement }) => (
  <div style={{ maxHeight: "300px", display: "flex" }}>
    <Table headings={["Column 1", "Column 2"]} {...props}>
      {[...Array(10)].map((_, index: number) => React.cloneElement(row, { key: index }))}
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
