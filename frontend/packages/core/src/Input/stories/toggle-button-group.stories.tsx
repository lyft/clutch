import * as React from "react";
import { ToggleButton } from "@mui/material";
import type { Meta } from "@storybook/react";

import ToggleButtonGroup from "../toggle-button-group";

export default {
  title: "Core/Input/ToggleButtonGroup",
  component: ToggleButtonGroup,
} as Meta;

/** Note this demo shows a usecase where an author wants to force a selection,
 * and swallow nulls / empty arrays that get passed when the user clicks
 * a button to deselect.
 */
export const EnforcedSelectionDemo = () => {
  const [value, setValue] = React.useState("MEOW");
  const onChange = (_: React.ChangeEvent<{}>, newValue: string) => {
    if (newValue) {
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

export const MultipleTrueVerticalDemo = () => {
  const [value, setValue] = React.useState<string[]>(["MEOW", "Ingress"]);
  const onChange = (_: React.ChangeEvent<{}>, newValue: string[]) => {
    setValue(newValue);
  };

  return (
    <ToggleButtonGroup
      multiple
      value={value}
      onChange={onChange}
      orientation="vertical"
      size="large"
    >
      <ToggleButton value="MEOW">MEOW</ToggleButton>
      <ToggleButton value="Ingress">Ingress</ToggleButton>
      <ToggleButton value="Egress">Egress</ToggleButton>
    </ToggleButtonGroup>
  );
};

export const SmallWithManyItemsDemo = () => {
  const [value, setValue] = React.useState("baaaaaa");
  const onChange = (_: React.ChangeEvent<{}>, newValue: string) => {
    setValue(newValue);
  };

  return (
    <ToggleButtonGroup value={value} onChange={onChange} size="small">
      <ToggleButton value="baaaaaa">baaaaaa</ToggleButton>
      <ToggleButton value="WOOF">WOOF</ToggleButton>
      <ToggleButton value="Bark">Bark</ToggleButton>
      <ToggleButton value="Chirp">Chirp</ToggleButton>
      <ToggleButton value="Moooooooooo">Moooooooooo</ToggleButton>
      <ToggleButton value="Squeek">Squeek</ToggleButton>
    </ToggleButtonGroup>
  );
};
