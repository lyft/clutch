import * as React from "react";
import styled from "@emotion/styled";
import type { SelectProps as MuiSelectProps } from "@material-ui/core";
import {
  FormControl as MuiFormControl,
  FormHelperText as MuiFormHelperText,
  InputLabel as MuiInputLabel,
  MenuItem,
  Select as MuiSelect,
} from "@material-ui/core";
import ErrorIcon from "@material-ui/icons/Error";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";

const StyledFormControl = styled(MuiFormControl)({
  "label + .MuiInput-formControl": {
    marginTop: "20px",
  },
});

const StyledFormHelperText = styled(MuiFormHelperText)({
  verticalAlign: "middle",
  display: "flex",
  position: "relative",
  fontSize: "12px",
  lineHeight: "16px",
  marginTop: "7px",
  color: "grey",

  "&.Mui-error": {
    color: "#db3615",
  },

  svg: {
    height: "16px",
    width: "16px",
    marginRight: "4px",
  },
});

const StyledInputLabel = styled(MuiInputLabel)({
  marginLeft: "2px",
  fontWeight: 500,
  fontSize: "14px",
  transform: "scale(1)",
  marginBottom: "10px",
  color: "rgba(13, 16, 48, 0.6)",
  "&.Mui-focused": {
    color: "rgba(13, 16, 48, 0.6)",
  },
  "&.Mui-error": {
    color: "#db3615",
  },
});

const SelectIcon = (props: any) => (
  <div {...props}>
    <ExpandMoreIcon />
  </div>
);

const BaseSelect = ({ ...props }: MuiSelectProps) => (
  <MuiSelect
    disableUnderline
    fullWidth
    IconComponent={SelectIcon}
    MenuProps={{
      style: { marginTop: "2px" },
      anchorOrigin: { vertical: "bottom", horizontal: "left" },
      transformOrigin: { vertical: "top", horizontal: "left" },
      getContentAnchorEl: null,
    }}
    {...props}
  />
);

const StyledSelect = styled(BaseSelect)({
  padding: "0",

  ".MuiSelect-root": {
    padding: "15px 60px 13px 16px",
    borderRadius: "4px",
    borderStyle: "solid",
    borderWidth: "1px",
    borderColor: "rgba(13, 16, 48, 0.38)",
    height: "20px",
  },

  "&.Mui-focused > .MuiSelect-root": {
    borderColor: "#3548d4",
  },

  "&.Mui-error > .MuiSelect-root": {
    borderColor: "#db3615",
  },

  ".MuiSelect-root.Mui-disabled": {
    backgroundColor: "rgba(13, 16, 48, 0.12)",
  },

  ".MuiSelect-icon": {
    height: "100%",
    width: "48px",
    top: "unset",
    transform: "unset",
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    boxSizing: "border-box",
  },

  ".MuiSelect-icon.Mui-disabled > svg": {
    color: "rgba(13, 16, 48, 0.6)",
  },

  ".MuiSelect-icon > svg": {
    color: "rgba(13, 16, 48, 0.6)",
    position: "absolute",
  },

  "&.Mui-focused > .MuiSelect-icon > svg": {
    color: "#0d1030",
  },

  ".MuiSelect-icon.MuiSelect-iconOpen > svg": {
    transform: "rotate(180deg)",
  },

  ".MuiSelect-select:focus": {
    backgroundColor: "inherit",
  },
});

interface SelectOption {
  label: string;
  value?: string;
}

export interface SelectProps extends Pick<MuiSelectProps, "disabled" | "error"> {
  defaultOption?: number;
  helperText?: string;
  label?: string;
  name: string;
  options: SelectOption[];
  onChange: (value: string) => void;
}

export const Select = ({
  defaultOption = 0,
  disabled,
  error,
  helperText,
  label,
  name,
  options,
  onChange,
}: SelectProps) => {
  if (options.length === 0) {
    return null;
  }

  const defaultIdx = defaultOption < options.length ? defaultOption : 0;
  const [selectedIdx, setSelectedIdx] = React.useState(defaultIdx);

  const updateSelectedOption = (event: React.ChangeEvent<{ name?: string; value: string }>) => {
    const { value } = event.target;
    const optionValues = options.map(o => o?.value || o.label);
    setSelectedIdx(optionValues.indexOf(value));
    onChange(value);
  };

  React.useEffect(() => {
    onChange(options[selectedIdx]?.value || options[selectedIdx].label);
  }, []);

  return (
    <StyledFormControl key={name} fullWidth disabled={disabled} error={error}>
      {label && <StyledInputLabel>{label}</StyledInputLabel>}
      <StyledSelect
        value={options[selectedIdx]?.value || options[selectedIdx].label}
        onChange={updateSelectedOption}
      >
        {options.map(option => (
          <MenuItem key={option.label} value={option?.value || option.label}>
            {option.label}
          </MenuItem>
        ))}
      </StyledSelect>
      {helperText && (
        <StyledFormHelperText>
          {error && <ErrorIcon />}
          {helperText}
        </StyledFormHelperText>
      )}
    </StyledFormControl>
  );
};

export default Select;
