import * as React from "react";
import { action } from "@storybook/addon-actions";
import type { Meta } from "@storybook/react";

import { WizardContext } from "../../Contexts";
import type { MetadataTableProps } from "../metadata-table";
import { MetadataTable } from "../metadata-table";

export default {
  title: "Core/Table/Metadata Table",
  decorators: [
    story => (
      <WizardContext.Provider
        value={() => {
          return {
            onSubmit: () => {},
            setOnSubmit: () => {},
            setIsLoading: () => {},
            displayWarnings: () => {},
            onBack: () => {},
          };
        }}
      >
        {story()}
      </WizardContext.Provider>
    ),
  ],
  component: MetadataTable,
} as Meta;

const Template = (props: MetadataTableProps) => (
  <div style={{ display: "flex", maxWidth: "800px" }}>
    <MetadataTable {...props} />
  </div>
);

export const Primary = Template.bind({});
Primary.args = {
  data: [
    { name: "Instance ID", value: "i-01cb1d09a5a1801e9" },
    { name: "Region", value: "us-east-1" },
    { name: "State", value: "RUNNING" },
    { name: "Instance Type", value: "c5.large" },
    { name: "Public IP Address", value: "54.234.102.49" },
    { name: "Private IP Address", value: "10.46.191.123" },
    { name: "Availability Zone", value: "us-east-1d" },
  ],
};

export const WithMutableRows = Template.bind({});
WithMutableRows.args = {
  onUpdate: action("update event"),
  data: [
    { name: "Name", value: "my-asg-name" },
    { name: "Region", value: "us-mock-1" },
    {
      name: "Min Size",
      value: 15,
      input: {
        type: "number",
        key: "size.min",
      },
    },
    {
      name: "Max Size",
      value: 25,
      input: { type: "number", key: "size.max" },
    },
    {
      name: "Desired Size",
      value: 20,
      input: { type: "number", key: "size.desired" },
    },
    { name: "Availability Zones", value: "us-mock-1b" },
  ],
};
