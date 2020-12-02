import * as React from "react";
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
