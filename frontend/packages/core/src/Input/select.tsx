import * as React from "react";
import {
  FormControl as MuiFormControl,
  FormHelperText as MuiFormHelperText,
  InputLabel as MuiInputLabel,
  MenuItem,
  Select as MuiSelect,
  SelectProps as MuiSelectProps,
} from "@material-ui/core";
import styled from "@emotion/styled";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import ErrorIcon from "@material-ui/icons/Error";

const StyledFormControl = styled(MuiFormControl)({
  "label + .MuiInput-formControl": {
    marginTop: "21px",
  },
});

const StyledFormHelperText = styled(MuiFormHelperText)({
  verticalAlign: "middle",
  display: "flex",
  position: "relative",
  fontSize: "12px",
  marginTop: "7px",
  color: "grey",

  "&.Mui-error": {
    color: "#db3615",
  },

  "svg": {
    height: "16px",
    width: "16px",
    marginRight: "5px",
  }
});

const StyledInputLabel = styled(MuiInputLabel)({
  fontWeight: "bold",
  fontSize: "13px",
  transform: "scale(1)",
  marginBottom: "10px",
  color: "rgba(13, 16, 48, 0.6)",
  "&.Mui-focused": {
    color: "rgba(13, 16, 48, 0.6)"
  },
  "&.Mui-error": {
    color: "#db3615",
  }
});

const SelectIcon = (props: any) => {
  return (
    <div {...props}>
      <ExpandMoreIcon />
    </div>
  );
};


const BaseSelect = ({ ...props }: MuiSelectProps) => (
  <MuiSelect
    disableUnderline
    fullWidth
    IconComponent={SelectIcon}
    MenuProps={{
      style: { marginTop: "2px" },
      anchorOrigin: { vertical: "bottom", horizontal: "left" },
      transformOrigin: { vertical: "top", horizontal: "left" },
      getContentAnchorEl: null
    }}
    {...props}
  />
)

const StyledSelect = styled(BaseSelect)({
  border: "1px solid rgba(13, 16, 48, 0.38)",
  borderRadius: "4px",
  padding: "0",
  marginTop: "6px",

  "&.Mui-focused": {
    borderColor: "#3548d4",
  },

  "&.Mui-error": {
    borderColor: "#db3615"
  },

  ".MuiSelect-root": {
    padding: "14px 16px",
  },

  ".MuiSelect-root.Mui-disabled": {
    backgroundColor: "rgba(13, 16, 48, 0.12)",
  },

  ".MuiSelect-icon": {
    borderLeftColor: "inherit",
    borderLeftWidth: "1px",
    borderLeftStyle: "solid",
    height: "100%",
    width: "48px",
    top: "unset",
    transform: "unset",
  },

  ".MuiSelect-icon.Mui-disabled > svg": {
    color: "rgba(13, 16, 48, 0.6)"
  },

  ".MuiSelect-icon > svg": {
    color: "rgba(13, 16, 48, 0.6)",
    position: "absolute",
    left: "50%",
    top: "50%",
    transform: "translate(-50%, -50%)"
  },

  "&.Mui-focused > .MuiSelect-icon > svg": {
    color: "#0d1030",
  },

  ".MuiSelect-icon.MuiSelect-iconOpen > svg": {
    transform: "translate(-50%, -50%) rotate(180deg)",

  },

  ".MuiSelect-select:focus": {
    backgroundColor: "inherit"
  }
})

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
      {helperText && <StyledFormHelperText>{error && <ErrorIcon />}{helperText}</StyledFormHelperText>}
    </StyledFormControl>
  );
};

export default Select;
