import * as React from "react";
import type { Meta } from "@storybook/react";

import type { TableProps } from "../table";
import { Table, TableRow } from "../table";

export default {
  title: "Core/Table/Table",
  component: Table,
} as Meta;

const Template = (props: TableProps) => (
  <div style={{ maxHeight: "300px", display: "flex" }}>
    <Table headings={["Column 1", "Column 2"]} {...props}>
      {[...Array(10)].map((_, index: number) => (
        // eslint-disable-next-line react/no-array-index-key
        <TableRow key={index}>
          <div>Value 1</div>
          <div>Value 2</div>
        </TableRow>
      ))}
    </Table>
  </div>
);

export const Primary = Template.bind({});

export const StickHeader = Template.bind({});
StickHeader.args = {
  stickyHeader: true,
};
