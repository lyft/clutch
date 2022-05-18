import * as React from "react";
import { ToggleButton } from "@material-ui/lab";
import type { Meta } from "@storybook/react";

import ToggleButtonGroup from "../toggle-button-group";

export default {
  title: "Core/Input/ToggleButtonGroup",
  component: ToggleButtonGroup,
} as Meta;

export const MultipleFalseDemo = () => {
  const [value, setValue] = React.useState("MEOW");
  const onChange = (_: React.ChangeEvent<{}>, newValue: string) => {
    // console.log(newValue);
    // if (newValue) {
    // console.log("newValue is not null" + newValue);
    // }
    setValue(newValue);
  };

  return (
    <ToggleButtonGroup multiple={false} value={value} onChange={onChange}>
      {/* console.log("value is " + value) */}
      <ToggleButton value="MEOW">MEOW</ToggleButton>
      <ToggleButton value="Ingress">Ingress</ToggleButton>
      <ToggleButton value="Egress">Egress</ToggleButton>
    </ToggleButtonGroup>
  );
};

export const MultipleTrueVerticalDemo = () => {
  const [value, setValue] = React.useState("MEOW");
  const onChange = (_: React.ChangeEvent<{}>, newValue: string) => {
    setValue(newValue);
  };

  return (
    <ToggleButtonGroup value={value} onChange={onChange} orientation="vertical" size="large">
      <ToggleButton value="MEOW">MEOW</ToggleButton>
      <ToggleButton value="Ingress">Ingress</ToggleButton>
      <ToggleButton value="Egress">Egress</ToggleButton>
    </ToggleButtonGroup>
  );
};

export const SmallWithManyItemsDemo = () => {
  const [value, setValue] = React.useState("MEOW");
  const onChange = (_: React.ChangeEvent<{}>, newValue: string) => {
    setValue(newValue);
  };

  return (
    <ToggleButtonGroup value={value} onChange={onChange} size="small">
      <ToggleButton value="MEOW">MEOW</ToggleButton>
      <ToggleButton value="WOOF">WOOF</ToggleButton>
      <ToggleButton value="Bark">Bark</ToggleButton>
      <ToggleButton value="Chirp">Chirp</ToggleButton>
      <ToggleButton value="Moooooooooo">Moooooooooo</ToggleButton>
      <ToggleButton value="Squeek">Squeek</ToggleButton>
    </ToggleButtonGroup>
  );
};
