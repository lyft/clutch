import React from "react";
import styled from "@emotion/styled";
import type { CheckboxProps as MuiCheckboxProps } from "@material-ui/core";
import {
  Checkbox as MuiCheckbox,
  FormControl as MuiFormControl,
  FormControlLabel,
  FormGroup,
  FormLabel,
  Grid,
} from "@material-ui/core";
import CheckIcon from "@material-ui/icons/Check";

const FormControl = styled(MuiFormControl)({
  width: "75%",
});

const StyledCheckbox = styled(MuiCheckbox)({
  color: "#6e7083",
  height: "48px",
  width: "48px",
  borderRadius: "30px",
  "&:hover": {
    background: "#f5f6fd",
  },
  "&:active": {
    background: "#d7daf6",
  },
  "&.Mui-checked": {
    color: "#ffffff",
    "&:hover": {
      background: "#f5f6fd",
    },
    "&:active": {
      background: "#d7daf6",
    },
    "&.Mui-disabled": {
      color: "#e7e7ea",
      ".MuiIconButton-label": {
        color: "rgba(13, 16, 48, 0.38)",
      },
    },
  },
});

const Icon = styled.div(
  {
    borderRadius: "2px",
    boxSizing: "border-box",
    height: "24px",
    width: "24px",
  },
  props => ({
    border: props["data-disabled"] ? "1px solid #e7e7ea" : "1px solid #6e7083",
  })
);

const SelectedIcon = styled.div(
  {
    borderRadius: "2px",
    boxSizing: "border-box",
    height: "24px",
    width: "24px",
    ".MuiSvgIcon-root": {
      display: "block",
    },
  },
  props => ({
    background: props["data-disabled"] ? "#e7e7eA" : "#3548d4",
  })
);

export interface CheckboxProps
  extends Pick<MuiCheckboxProps, "checked" | "disabled" | "name" | "onChange"> {}

const Checkbox: React.FC<CheckboxProps> = ({ checked, disabled, ...props }) => {
  return (
    <StyledCheckbox
      checked={checked}
      icon={<Icon data-disabled={disabled} />}
      checkedIcon={
        <SelectedIcon data-disabled={disabled}>
          <CheckIcon />
        </SelectedIcon>
      }
      {...props}
      disabled={disabled}
    />
  );
};

export interface CheckboxPanelProps {
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

export { CheckboxPanel, Checkbox };
