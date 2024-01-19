import * as React from "react";
import CheckIcon from "@mui/icons-material/Check";
import type { CheckboxProps as MuiCheckboxProps, Theme } from "@mui/material";
import {
  alpha,
  Checkbox as MuiCheckbox,
  FormControl as MuiFormControl,
  FormControlLabel,
  FormGroup,
  FormLabel,
  Grid,
} from "@mui/material";

import styled from "../styled";

const FormControl = styled(MuiFormControl)({
  width: "75%",
});

const StyledCheckbox = styled(MuiCheckbox)(({ theme }: { theme: Theme }) => ({
  color: theme.palette.secondary[400],
  borderRadius: "50%",
  "&:hover": {
    background: theme.palette.primary[100],
  },
  "&:active": {
    background: theme.palette.primary[300],
  },
  "&.Mui-checked": {
    color: theme.palette.contrastColor,
    "&:hover": {
      background: theme.palette.primary[100],
    },
    "&:active": {
      background: theme.palette.primary[300],
    },
    "&.Mui-disabled": {
      color: theme.palette.secondary[200],
      ".MuiIconButton-label": {
        color: alpha(theme.palette.secondary[900], 0.38),
      },
    },
  },
}));

type Size = "20px" | "24px";

interface StyledIconProps {
  $disabled: CheckboxProps["disabled"];
  $size: Size;
}

const Icon = styled("div")<StyledIconProps>(
  {
    borderRadius: "2px",
    boxSizing: "border-box",
  },
  props => ({ theme }: { theme: Theme }) => ({
    height: props.$size,
    width: props.$size,
    border: props.$disabled
      ? `1px solid ${theme.palette.secondary[200]}`
      : `1px solid ${theme.palette.secondary[400]}`,
  })
);

const SelectedIcon = styled("div")<StyledIconProps>(
  {
    borderRadius: "2px",
    boxSizing: "border-box",
    ".MuiSvgIcon-root": {
      display: "block",
    },
  },
  props => ({ theme }: { theme: Theme }) => ({
    height: props.$size,
    width: props.$size,
    background: props.$disabled ? theme.palette.secondary[200] : theme.palette.primary[600],
    ".MuiSvgIcon-root": {
      height: props.$size,
      width: props.$size,
    },
  })
);

export interface CheckboxProps
  extends Pick<MuiCheckboxProps, "checked" | "disabled" | "name" | "onChange" | "size"> {}

// TODO (sperry): add 16px size variant
const Checkbox: React.FC<CheckboxProps> = ({ checked, disabled = false, size, ...props }) => {
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
      icon={<Icon $disabled={disabled} $size={sizePx} />}
      checkedIcon={
        <SelectedIcon $disabled={disabled} $size={sizePx}>
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
