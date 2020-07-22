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
  header: string;
  options: {
    [option: string]: string;
  };
  onChange: (option: string, checked: boolean) => null;
}

const CheckboxPanel: React.FC<CheckboxPanelProps> = ({ header, options, onChange }) => {
  const allOptions = {};
  Object.keys(options).forEach(option => {
    allOptions[option] = { checked: false, value: options[option] };
  });
  const [selected, setSelected] = React.useState(allOptions);

  const onToggle = e => {
    const selectedOption = { ...selected[e.target.name], checked: e.target.checked };
    const updatedSelections = { ...selected, [e.target.name]: selectedOption };
    setSelected(updatedSelections);
    onChange(allOptions[e.target.name].value, e.target.checked);
  };

  return (
    <FormControl>
      <Grid container direction="column">
        <FormLabel color="secondary" component="legend">
          {header}
        </FormLabel>
        <FormGroup row>
          {Object.keys(allOptions).map(option => (
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
