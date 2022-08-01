import * as React from "react";
import styled from "@emotion/styled";
import ErrorIcon from "@mui/icons-material/Error";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import type { SelectProps as MuiSelectProps } from "@mui/material";
import {
  FormControl as MuiFormControl,
  FormHelperText as MuiFormHelperText,
  InputLabel as MuiInputLabel,
  ListSubheader,
  MenuItem,
  Select as MuiSelect,
} from "@mui/material";
import { flatten } from "lodash";

const StyledFormControl = styled(MuiFormControl)({
  margin: "8px 8px 8px 0px",
});

const StyledFormHelperText = styled(MuiFormHelperText)({
  alignItems: "center",
  display: "flex",
  position: "relative",
  fontSize: "12px",
  lineHeight: "16px",
  marginTop: "7px",
  marginLeft: "0px",
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
  "--label-default-color": "rgba(13, 16, 48, 0.6)",

  color: "var(--label-default-color)",
  "&.Mui-focused": {
    color: "var(--label-default-color)",
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

const BaseSelect = ({ className, ...props }: MuiSelectProps) => (
  <MuiSelect
    // disableUnderline
    fullWidth
    IconComponent={SelectIcon}
    className={className}
    MenuProps={{
      classes: {
        list: className,
        paper: className,
      },
      anchorOrigin: { vertical: "bottom", horizontal: "left" },
      transformOrigin: { vertical: "top", horizontal: "left" },
    }}
    {...props}
  />
);

const StyledSelect = styled(BaseSelect)({
  "--notched-border-width": "1px",
  padding: "0",
  backgroundColor: "#FFFFFF",
  height: "unset",

  ".MuiOutlinedInput-notchedOutline": {
    borderColor: "rgba(13, 16, 48, 0.38)",
    borderWidth: "var(--notched-border-width)",
  },

  "&.Mui-focused": {
    "> .MuiSelect-icon > svg": {
      color: "#0d1030",
    },
    "> .MuiOutlinedInput-notchedOutline": {
      borderColor: "#3458d4",
      borderWidth: "var(--notched-border-width)",
    },
  },
  "&.Mui-error > .MuiOutlinedInput-notchedOutline": {
    borderColor: "#db3615",
    borderWidth: "var(--notched-border-width)",
  },

  ".MuiSelect-select": {
    height: "20px",
    display: "flex",
    padding: "15px 60px 13px 16px",

    ":focus": {
      backgroundColor: "inherit",
    },

    "&.Mui-disabled": {
      backgroundColor: "rgba(13, 16, 48, 0.12)",
    },
  },

  ul: {
    borderRadius: "4px",
    border: "1px solid rgba(13, 16, 48, 0.1)",
  },

  ".MuiMenuItem-root": {
    color: "#0d1030",
    height: "48px",

    ":first-of-type": {
      borderRadius: "4px 4px 0 0",
    },
    ":last-child": {},
    ":hover": {
      backgroundColor: "#e7e7ea",
    },
    ":active": {
      backgroundColor: "#dbdbe0",
    },
    "&.Mui-selected": {
      backgroundColor: "rgba(53, 72, 212, 0.1)",
      ":hover": {
        backgroundColor: "rgba(53, 72, 212, 0.1)",
      },
    },
  },

  "&.MuiMenu-paper": {
    marginTop: "5px",
    border: "none",
    boxShadow: "0px 5px 15px rgba(53, 72, 212, 0.2)",
    maxHeight: "50vh",
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

    "> svg": {
      color: "rgba(13, 16, 48, 0.6)",
      position: "absolute",
    },

    "&.MuiSelect-iconOpen > svg": {
      transform: "rotate(180deg)",
    },

    "&.Mui-disabled > svg": {
      color: "rgba(13, 16, 48, 0.6)",
    },
  },

  ".MuiListSubheader-root": {
    color: "#939495",
    cursor: "default",
    pointerEvents: "none", // disables the select from closing on clicking the subheader
  },
});

interface BaseSelectOptions {
  label: string;
  value?: string;
  startAdornment?: React.ReactElement;
}

export interface SelectOption extends BaseSelectOptions {
  group?: BaseSelectOptions[];
}

export interface SelectProps extends Pick<MuiSelectProps, "disabled" | "error"> {
  defaultOption?: number;
  helperText?: string;
  label?: string;
  name: string;
  options: SelectOption[];
  onChange?: (value: string) => void;
}

const Select = ({
  defaultOption = 0,
  disabled,
  error,
  helperText,
  label,
  name,
  options,
  onChange,
}: SelectProps) => {
  const defaultIdx = defaultOption < options.length && defaultOption > 0 ? defaultOption : 0;
  const [selectedIdx, setSelectedIdx] = React.useState(defaultIdx);

  // Flattens all options and sub grouped options for easier retrieval
  const flatOptions: BaseSelectOptions[] = flatten(
    options.map((option: SelectOption) =>
      option.group ? option.group.map(groupOption => groupOption) : option
    )
  );

  React.useEffect(() => {
    if (flatOptions.length !== 0) {
      onChange && onChange(flatOptions[selectedIdx]?.value || flatOptions[selectedIdx].label);
    }
  }, []);

  const updateSelectedOption = event => {
    const { value } = event.target;
    // handle if selecting a header option
    if (!value) {
      return;
    }
    setSelectedIdx(flatOptions.findIndex(opt => opt?.value === value || opt?.label === value));
    onChange && onChange(value);
  };

  if (flatOptions.length === 0) {
    return null;
  }

  const menuItem = option => (
    <MenuItem key={option.label} value={option?.value || option.label}>
      {option?.startAdornment &&
        React.cloneElement(option.startAdornment, {
          style: { height: "100%", marginRight: "8px", ...option.startAdornment.props.style },
        })}
      {option.label}
    </MenuItem>
  );

  const renderSelectItems = option => {
    if (option.group) {
      return [
        <ListSubheader>{option.label}</ListSubheader>,
        option.group.map(opt => menuItem(opt)),
      ];
    }
    return menuItem(option);
  };

  return (
    <StyledFormControl id={name} key={name} fullWidth disabled={disabled} error={error}>
      {label && <StyledInputLabel>{label}</StyledInputLabel>}
      {flatOptions.length && (
        <StyledSelect
          id={`${name}-select`}
          value={flatOptions[selectedIdx]?.value || flatOptions[selectedIdx].label}
          onChange={updateSelectedOption}
          label={label}
        >
          {options?.map(option => renderSelectItems(option))}
        </StyledSelect>
      )}
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
