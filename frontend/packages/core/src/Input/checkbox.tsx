import React from "react";
import styled from "@emotion/styled";
import type { CheckboxProps } from "@material-ui/core";
import {
  Checkbox as MuiCheckbox,
  FormControl as MuiFormControl,
  FormControlLabel,
  FormGroup,
  FormLabel,
  Grid,
} from "@material-ui/core";

const FormControl = styled(MuiFormControl)`
  width: 75%;
`;

const Checkbox = styled(MuiCheckbox)`
  &.MuiCheckbox-root {
    color: #0d1030;
    &:hover {
      color: #3548d4;
      background-color: transparent;
    }
  }

  &.Mui-checked {
    color: #3548d4;
    &:hover {
      background-color: transparent;
    }
  }

  &.Mui-disabled {
    color: #0d1030;
    opacity: 0.38;
  }
`;

export interface CheckboxPanelProps extends Pick<CheckboxProps, "disabled"> {
  header?: string;
  options: {
    [option: string]: boolean;
  };
  onChange: (state: { [option: string]: boolean }) => void;
}

const CheckboxPanel: React.FC<CheckboxPanelProps> = ({ header, options, onChange, ...props }) => {
  const allOptions = {};
  Object.keys(options).forEach(option => {
    allOptions[option] = { checked: options[option], value: option };
  });
  const [selected, setSelected] = React.useState(allOptions);

  const onToggle = e => {
    const targetName = e.target.name;
    const targetValue = e.target.checked;
    const selectedOption = { ...selected[targetName], checked: targetValue };
    const updatedSelections = { ...selected, [targetName]: selectedOption };
    setSelected(updatedSelections);
    const callbackOptions = {};
    Object.keys(allOptions).forEach(option => {
      callbackOptions[option] =
        option === targetName
          ? targetValue
          : selected[option]
          ? selected[option].checked
          : allOptions[option].checked;
    });
    onChange(callbackOptions);
  };

  const optionKeys = Object.keys(allOptions);
  const column1Keys = [...optionKeys].splice(0, Math.ceil(optionKeys.length / 2));
  const column2Keys = [...optionKeys].splice(column1Keys.length, optionKeys.length);

  return (
    <FormControl>
      <Grid container direction="column">
        <FormLabel color="secondary" focused>
          {header}
        </FormLabel>
        <Grid container direction="row">
          <FormGroup>
            {column1Keys.map(option => (
              <FormGroup row key={option}>
                <FormControlLabel
                  key={option}
                  control={
                    <Checkbox
                      checked={selected[option].checked}
                      onChange={onToggle}
                      name={option}
                      {...props}
                    />
                  }
                  label={option}
                />
              </FormGroup>
            ))}
          </FormGroup>
          <FormGroup>
            {column2Keys.map(option => (
              <FormGroup row key={option}>
                <FormControlLabel
                  key={option}
                  control={
                    <Checkbox
                      checked={selected[option].checked}
                      onChange={onToggle}
                      name={option}
                      {...props}
                    />
                  }
                  label={option}
                />
              </FormGroup>
            ))}
          </FormGroup>
        </Grid>
      </Grid>
    </FormControl>
  );
};

export default CheckboxPanel;
