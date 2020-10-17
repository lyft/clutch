import React from "react";
import { FormControl, FormControlLabel, FormLabel, Radio, RadioGroup } from "@material-ui/core";

interface RadioControlItem {
  label: string;
  value: string;
}

interface RadioControlProps {
  name: string;
  label: string;
  items: RadioControlItem[];
  onChange: (value: string) => void;
}

const RadioControl: React.FC<RadioControlProps> = ({ name, label, items, onChange }) => {
  return (
    <FormControl key={name}>
      <FormLabel component="legend">Upstream Cluster Type</FormLabel>
      <RadioGroup
        aria-label={label}
        name={name}
        defaultValue={items[0].value}
        onChange={e => onChange(e.target.value)}
      >
        {items &&
          items.map(item => {
            return (
              <FormControlLabel
                key={item.value}
                value={item.value}
                control={<Radio />}
                label={item.label}
              />
            );
          })}
      </RadioGroup>
    </FormControl>
  );
};

export { RadioControl };
