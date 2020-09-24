import React from "react";
import {
  Checkbox as MuiCheckbox,
  FormControl as MuiFormControl,
  FormControlLabel,
  FormGroup,
  FormLabel,
  Grid,
} from "@material-ui/core";
import styled from "styled-components";

const FormControl = styled(MuiFormControl)`
  width: 75%;
`;

const Checkbox = styled(MuiCheckbox)`
  ${({ theme }) => `
  color: ${theme.palette.text.primary};
  `}
`;

interface CheckboxPanelProps {
  header?: string;
  options: {
    [option: string]: boolean;
  };
  onChange: (state: { [option: string]: boolean }) => void;
}

const CheckboxPanel: React.FC<CheckboxPanelProps> = ({ header, options, onChange }) => {
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
  const column1Keys = [...optionKeys].splice(0, optionKeys.length / 2);
  const column2Keys = [...optionKeys].splice(column1Keys.length + 1, optionKeys.length);

  return (
    <FormControl>
      <Grid container direction="row" justify="center" alignItems="stretch">
        <FormLabel color="secondary" component="legend">
          {header}
        </FormLabel>
        <FormGroup>
          {column1Keys.map(option => (
            <FormGroup row key={option}>
              <FormControlLabel
                key={option}
                control={
                  <Checkbox checked={selected[option].checked} onChange={onToggle} name={option} />
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
                  <Checkbox checked={selected[option].checked} onChange={onToggle} name={option} />
                }
                label={option}
              />
            </FormGroup>
          ))}
        </FormGroup>
      </Grid>
    </FormControl>
  );
};

export default CheckboxPanel;
