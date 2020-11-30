import React from "react";
import {
  FormControl as MuiFormControl,
  InputLabel as MuiInputLabel,
  MenuItem,
  Select as MuiSelect,
  SelectProps as MuiSelectProps,
} from "@material-ui/core";
import styled from "@emotion/styled";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";


const StyledFormControl = styled(MuiFormControl)({
  "label + .MuiInput-formControl": {
    marginTop: "21px",
  }
});

const StyledInputLabel = styled(MuiInputLabel)({
  fontWeight: "bold",
  fontSize: "13px",
  transform: "scale(1)",
  marginBottom: "10px",
  color: "grey",
  "&.Mui-focused": {
    "color": "grey"
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

  ".MuiSelect-root": {
    padding: "14px 16px"
  },

  ".MuiSelect-icon": {
    height: "100%",
    width: "50px",
    top: "unset",
    borderLeft: "1px solid rgba(13, 16, 48, 0.38)",
    transform: "unset"
  },

  ".MuiSelect-icon > svg": {
    color: "black",
    position: "absolute",
    left: "50%",
    top: "50%",
    transform: "translate(-50%, -50%)"
  },

  ".MuiSelect-icon.MuiSelect-iconOpen > svg": {
    transform: "translate(-50%, -50%) rotate(180deg)"
  },

  ".MuiSelect-select:focus": {
    backgroundColor: "inherit"
  }
})

interface SelectOption {
  label: string;
  value?: string;
}

export interface SelectProps {
  defaultOption?: number;
  label?: string;
  maxWidth?: string;
  name: string;
  options: SelectOption[];
  onChange: (value: string) => void;
}

const Select: React.FC<SelectProps> = ({
  defaultOption = 0,
  label,
  name,
  options,
  onChange,
}) => {
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
    <StyledFormControl key={name} fullWidth>
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
    </StyledFormControl>
  );
};

export default Select;
