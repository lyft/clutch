import React from "react";
import {
  FormControl as MuiFormControl,
  InputLabel as MuiInputLabel,
  MenuItem,
  Select as MuiSelect,
} from "@material-ui/core";
import styled from "styled-components";

const InputLabel = styled(MuiInputLabel)`
  ${({ theme }) => `
  && {
    color: ${theme.palette.text.primary};
  }
  `}
`;

const FormControl = styled(MuiFormControl)`
  ${({ theme, ...props }) => `
  display: flex;
  width: 100%;
  max-width: ${props["data-max-width"] || "500px"};
  .MuiInput-underline:after {
    border-bottom: 2px solid ${theme.palette.accent.main};
  }
  `}
`;

const StyledSelect = styled(MuiSelect)`
  ${({ ...props }) => `
  display: flex;
  margin: 16px 0 16px 0;
  width: 100%;
  max-width: ${props["data-max-width"] || "500px"};
  `}
`;

interface SelectOption {
  label: string;
  value?: string;
}

interface SelectProps {
  defaultOption?: number;
  label: string;
  maxWidth?: string;
  name: string;
  options: SelectOption[];
  onChange: (value: string) => void;
}

const Select: React.FC<SelectProps> = ({
  defaultOption = 0,
  label,
  maxWidth,
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
    const optionValues = options.map(o => o.value || o.label);
    setSelectedIdx(optionValues.indexOf(value));
    onChange(value);
  };

  React.useEffect(() => {
    onChange(options[selectedIdx].value || options[selectedIdx].label);
  }, []);

  return (
    <FormControl key={name} data-max-width={maxWidth}>
      <InputLabel color="secondary">{label}</InputLabel>
      <StyledSelect
        value={options[selectedIdx].value || options[selectedIdx].label}
        onChange={updateSelectedOption}
      >
        {options.map(option => (
          <MenuItem key={option.label} value={option.value || option.label}>
            {option.label}
          </MenuItem>
        ))}
      </StyledSelect>
    </FormControl>
  );
};

export default Select;
