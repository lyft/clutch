import * as React from "react";
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
  borderRadius: "50%",
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

const Icon = styled.div<{ disabled: boolean; size: string }>(
  {
    borderRadius: "2px",
    boxSizing: "border-box",
  },
  props => ({
    height: props.size,
    width: props.size,
    border: props.disabled ? "1px solid #e7e7ea" : "1px solid #6e7083",
  })
);

const SelectedIcon = styled.div<{ disabled: boolean; size: string }>(
  {
    borderRadius: "2px",
    boxSizing: "border-box",
    ".MuiSvgIcon-root": {
      display: "block",
    },
  },
  props => ({
    height: props.size,
    width: props.size,
    background: props.disabled ? "#e7e7eA" : "#3548d4",
    ".MuiSvgIcon-root": {
      height: props.size,
      width: props.size,
    },
  })
);

export interface CheckboxProps
  extends Pick<MuiCheckboxProps, "checked" | "disabled" | "name" | "onChange" | "size"> {}

// TODO (sperry): add 16px size variant
const Checkbox: React.FC<CheckboxProps> = ({ checked, disabled, size, ...props }) => {
  let sizePx;
  switch (size) {
    case "small":
      sizePx = "20px";
      break;
    default:
      sizePx = "24px";
  }

  return (
    <StyledCheckbox
      checked={checked}
      size={size}
      icon={<Icon disabled={disabled} size={sizePx} />}
      checkedIcon={
        <SelectedIcon disabled={disabled} size={sizePx}>
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
