import * as React from "react";
import type { Meta } from "@storybook/react";

import ToggleButtonGroup from "../toggle-button-group";

export default {
  title: "Core/Input/ToggleButtonGroup",
  component: ToggleButtonGroup,
} as Meta;

export const ExclusiveDemo = () => {
  const [value, setValue] = React.useState("meow");
  const onChange = (_event: React.ChangeEvent<{}>, newValue: string) => {
    // Note that we check for null because we want to enforce that
    // at least value is active
    if (newValue !== null) {
      setValue(newValue);
    }
  };

  return (
    <ToggleButtonGroup
      exclusive
      currentValue={value}
      onChange={onChange}
      toggleButtonValues={["Egress", "Ingress", "meow"]}
    />
  );
};
