import * as React from "react";
import type { Meta } from "@storybook/react";

import type { TabsProps } from "../tab";
import { Tab, Tabs } from "../tab";

export default {
  title: "Core/Tab/Tabs",
  component: Tab,
  argTypes: {
    onClick: { action: "onClick event" },
  },
} as Meta;

const Template = ({ tabCount, value, ...props }: TabsProps & { tabCount: number }) => {
  const [selectedValue, setSelectedValue] = React.useState(value - 1);
  return (
    <Tabs value={selectedValue} onChange={(_, v) => setSelectedValue(v)} {...props}>
      {[...Array(tabCount)].map((_, index: number) => (
        // eslint-disable-next-line react/no-array-index-key
        <Tab key={index} label={`Tab ${index + 1}`} value={index} />
      ))}
    </Tabs>
  );
};

export const Primary = Template.bind({});
Primary.args = {
  tabCount: 2,
  value: 1,
};
