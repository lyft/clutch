import * as React from "react";
import styled from "@emotion/styled";
import CancelIcon from "@mui/icons-material/Cancel";
import ErrorIcon from "@mui/icons-material/Error";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import type { SelectProps as MuiSelectProps, Theme } from "@mui/material";
import {
  alpha,
  FormControl as MuiFormControl,
  FormHelperText as MuiFormHelperText,
  InputLabel as MuiInputLabel,
  ListSubheader,
  MenuItem,
  Select as MuiSelect,
} from "@mui/material";
import { flatten } from "lodash";

import { Chip } from "../chip";

const StyledFormHelperText = styled(MuiFormHelperText)(({ theme }: { theme: Theme }) => ({
  alignItems: "center",
  display: "flex",
  position: "relative",
  fontSize: "12px",
  lineHeight: "16px",
  marginTop: "7px",
  marginLeft: "0px",
  color: "grey",

  "&.Mui-error": {
    color: theme.palette.error[600],
  },

  svg: {
    height: "16px",
    width: "16px",
    marginRight: "4px",
  },
}));

const StyledInputLabel = styled(MuiInputLabel)(({ theme }: { theme: Theme }) => ({
  "--label-default-color": alpha(theme.palette.secondary[900], 0.6),

  color: "var(--label-default-color)",
  "&.Mui-focused": {
    color: "var(--label-default-color)",
  },
  "&.Mui-error": {
    color: theme.palette.error[600],
  },
}));

const SelectIcon = (props: any) => (
  <div {...props}>
    <ExpandMoreIcon />
  </div>
);

const BaseSelect = ({ className, ...props }: MuiSelectProps) => (
  <MuiSelect
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

const StyledSelect = styled(BaseSelect)(({ theme }: { theme: Theme }) => ({
  "--notched-border-width": "1px",
  padding: "0",
  backgroundColor: theme.palette.contrastColor,
  minWidth: "fit-content",

  ".MuiOutlinedInput-notchedOutline": {
    borderColor: alpha(theme.palette.secondary[900], 0.38),
    borderWidth: "var(--notched-border-width)",
  },

  "&.Mui-focused": {
    "> .MuiSelect-icon > svg": {
      color: theme.palette.secondary[900],
    },
    "> .MuiOutlinedInput-notchedOutline": {
      borderColor: theme.palette.primary[600],
      borderWidth: "var(--notched-border-width)",
    },
  },
  "&.Mui-error > .MuiOutlinedInput-notchedOutline": {
    borderColor: theme.palette.error[600],
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
      backgroundColor: alpha(theme.palette.secondary[900], 0.12),
    },
  },

  ul: {
    borderRadius: "4px",
    border: `1px solid ${alpha(theme.palette.secondary[900], 0.1)}`,
  },

  ".MuiMenuItem-root": {
    color: theme.palette.secondary[900],
    height: "48px",

    ":first-of-type": {
      borderRadius: "4px 4px 0 0",
    },
    ":last-child": {},
    ":hover": {
      backgroundColor: theme.palette.secondary[200],
    },
    ":active": {
      backgroundColor: theme.palette.secondary[200],
    },
    "&.Mui-selected": {
      backgroundColor: alpha(theme.palette.primary[600], 0.1),
      ":hover": {
        backgroundColor: alpha(theme.palette.primary[600], 0.1),
      },
    },
  },

  "&.MuiMenu-paper": {
    marginTop: "5px",
    border: "none",
    boxShadow: `0px 5px 15px ${alpha(theme.palette.primary[600], 0.2)}`,
    maxHeight: "25vh",
  },

  ".MuiSelect-icon": {
    height: "100%",
    width: "48px",
    top: "unset",
    transform: "unset",
    display: "flex",
    justifyContent: "flex-end",
    alignItems: "center",
    boxSizing: "border-box",

    "> svg": {
      color: alpha(theme.palette.secondary[900], 0.6),
      position: "absolute",
    },

    "&.MuiSelect-iconOpen > svg": {
      transform: "rotate(180deg)",
    },

    "&.Mui-disabled > svg": {
      color: alpha(theme.palette.secondary[900], 0.6),
    },
  },

  ".MuiListSubheader-root": {
    color: theme.palette.secondary[400],
    cursor: "default",
    pointerEvents: "none", // disables the select from closing on clicking the subheader
  },
}));

interface BaseSelectOptions {
  label: string;
  value?: string;
  startAdornment?: React.ReactElement;
}

export interface SelectOption extends BaseSelectOptions {
  group?: BaseSelectOptions[];
}

const flattenBaseSelectOptions = (options: BaseSelectOptions[]) =>
  flatten(
    options.map((option: SelectOption) =>
      option.group ? option.group.map(groupOption => groupOption) : option
    )
  );

const menuItemFromOption = (option: BaseSelectOptions) => (
  <MenuItem key={option.label} value={option?.value || option.label}>
    {option?.startAdornment &&
      React.cloneElement(option.startAdornment, {
        style: {
          height: "100%",
          maxHeight: "20px",
          marginRight: "8px",
          ...option.startAdornment.props.style,
        },
      })}
    {option.label}
  </MenuItem>
);

const renderSelectItems = (option: SelectOption) => {
  if (option.group) {
    return [
      <ListSubheader>{option.label}</ListSubheader>,
      option.group.map(opt => menuItemFromOption(opt)),
    ];
  }
  return menuItemFromOption(option);
};

// Will take an array of strings or integers and attempt to find the indexes where they exist based on the flattened items
const calculateDefaultOptions = (
  defaultOptions: Array<number> | Array<string>,
  flattenedOptions: BaseSelectOptions[]
): Array<number> => {
  const options = [];

  if (defaultOptions === undefined || defaultOptions.length === 0) {
    return options;
  }

  defaultOptions.forEach(option => {
    if (Number.isInteger(option)) {
      options.push(option < flattenedOptions.length && option > 0 ? option : 0);
    }

    // we're a string, lets look it up based on the value/label and default to 0 if none
    const index = flattenedOptions?.findIndex(
      opt => opt?.value === option || opt?.label === option
    );
    options.push(index >= 0 ? index : 0);
  });
  return options;
};

export interface SelectProps extends Pick<MuiSelectProps, "disabled" | "error" | "value"> {
  defaultOption?: number | string;
  helperText?: string;
  label?: string;
  name: string;
  options: SelectOption[];
  onChange?: (value: string) => void;
  noDefault?: boolean;
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
  value,
  noDefault,
}: SelectProps) => {
  // Flattens all options and sub grouped options for easier retrieval
  const flatOptions: BaseSelectOptions[] = flattenBaseSelectOptions(options);
  const defaultOptions = calculateDefaultOptions(
    [defaultOption] as Array<number> | Array<string>,
    flatOptions
  );

  const [selectedIdx, setSelectedIdx] = React.useState(
    defaultOptions.length > 0 ? defaultOptions[0] : 0
  );

  React.useEffect(() => {
    if (flatOptions.length !== 0 && !noDefault) {
      onChange && onChange(flatOptions[selectedIdx]?.value || flatOptions[selectedIdx].label);
    }
  }, []);

  const updateSelectedOption = event => {
    const targetValue = event.target.value;

    // handle if selecting a header option
    if (!targetValue) {
      return;
    }

    setSelectedIdx(
      flatOptions.findIndex(opt => opt?.value === targetValue || opt?.label === targetValue)
    );

    onChange && onChange(targetValue);
  };

  if (flatOptions.length === 0) {
    return null;
  }

  return (
    <MuiFormControl id={name} key={name} disabled={disabled} error={error} fullWidth>
      {label && <StyledInputLabel>{label}</StyledInputLabel>}
      {flatOptions.length && (
        <StyledSelect
          id={`${name}-select`}
          value={value ?? (flatOptions[selectedIdx]?.value || flatOptions[selectedIdx].label)}
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
    </MuiFormControl>
  );
};

export interface MultiSelectProps extends Pick<MuiSelectProps, "disabled" | "error"> {
  defaultOptions?: Array<number> | Array<string>;
  helperText?: string;
  label?: string;
  name: string;
  selectOptions: SelectOption[];
  chipDisplay?: boolean;
  onChange?: (values: Array<string>) => void;
}

const MultiSelect = ({
  defaultOptions = [],
  disabled,
  error,
  helperText,
  label,
  name,
  selectOptions,
  chipDisplay = false,
  onChange,
}: MultiSelectProps) => {
  // Flattens all options and sub grouped options for easier retrieval
  const flatOptions: BaseSelectOptions[] = flattenBaseSelectOptions(selectOptions);

  const [selectedOptions, setSelectedOptions] = React.useState<Array<number>>(
    calculateDefaultOptions(defaultOptions, flatOptions)
  );

  const selectedValues = () =>
    selectedOptions.map(idx => flatOptions[idx].value || flatOptions[idx].label);

  React.useEffect(() => {
    if (flatOptions.length !== 0) {
      onChange && onChange(selectedValues());
    }
  }, [selectedOptions]);

  const updateSelectedOptions = event => {
    const { value } = event.target;
    // handle if selecting a header option
    if (!value) {
      return;
    }
    const findIndex = val => flatOptions.findIndex(opt => opt.value === val || opt.label === val);
    setSelectedOptions(value.map(val => findIndex(val)));
  };

  const onDeleteChip = value => () => {
    const findIndex = val => flatOptions.findIndex(opt => opt.value === val || opt.label === val);
    const updatedOptions = selectedOptions.filter(option => option !== findIndex(value));
    setSelectedOptions(updatedOptions);
  };

  if (flatOptions.length === 0) {
    return null;
  }

  return (
    <MuiFormControl id={name} key={name} disabled={disabled} error={error} fullWidth>
      {label && <StyledInputLabel>{label}</StyledInputLabel>}
      {flatOptions.length && (
        <StyledSelect
          multiple
          id={`${name}-multi-select`}
          value={selectedValues()}
          onChange={updateSelectedOptions}
          label={label}
          {...(chipDisplay && {
            renderValue: (selected: string[]) => (
              <div style={{ display: "flex", gap: "4px" }}>
                {selected.sort().map(value => (
                  <Chip
                    variant="neutral"
                    label={value}
                    key={value}
                    size="small"
                    onDelete={onDeleteChip(value)}
                    deleteIcon={<CancelIcon onMouseDown={event => event.stopPropagation()} />}
                  />
                ))}
              </div>
            ),
          })}
        >
          {selectOptions?.map(option => renderSelectItems(option))}
        </StyledSelect>
      )}
      {helperText && (
        <StyledFormHelperText>
          {error && <ErrorIcon />}
          {helperText}
        </StyledFormHelperText>
      )}
    </MuiFormControl>
  );
};

export { Select, MultiSelect };
