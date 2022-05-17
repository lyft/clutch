import * as React from "react";
import { ToggleButton } from "@material-ui/lab";
import type { Meta } from "@storybook/react";

import ToggleButtonGroup from "../toggle-button-group";

export default {
  title: "Core/Input/ToggleButtonGroup",
  component: ToggleButtonGroup,
} as Meta;

export const ExclusiveDemo = () => {
  const [value, setValue] = React.useState("MEOW");
  const onChange = (_: React.ChangeEvent<{}>, newValue: string) => {
    // Note that we check for null because we want to enforce that
    // at least value is active
    if (newValue !== null) {
      setValue(newValue);
    }
  };

  return (
    <ToggleButtonGroup multiple={false} value={value} onChange={onChange}>
      <ToggleButton value="MEOW">MEOW</ToggleButton>
      <ToggleButton value="Ingress">Ingress</ToggleButton>
      <ToggleButton value="Egress">Egress</ToggleButton>
    </ToggleButtonGroup>
  );
};
