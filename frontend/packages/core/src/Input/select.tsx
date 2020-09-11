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
  width: 100%;
  max-width: ${props["data-max-width"] || "500px"};
  `}
`;

interface SelectProps {
  defaultOption?: string;
  label: string;
  maxWidth?: string;
  name: string;
  options: string[];
  onChange: (value: string) => void;
}

const Select: React.FC<SelectProps> = ({
  defaultOption,
  label,
  maxWidth,
  name,
  options,
  onChange,
}) => {
  if (options.length === 0) {
    return null;
  }

  const defaultIdx = (options || []).indexOf(defaultOption);
  const [selectedValue, setSelectedValue] = React.useState(options[defaultIdx] || options[0]);

  const updateSelectedOption = (event: React.ChangeEvent<{ name?: string; value: string }>) => {
    setSelectedValue(event.target.value);
    onChange(event.target.value);
  };

  return (
    <FormControl key={name} data-max-width={maxWidth}>
      <InputLabel color="secondary">{label}</InputLabel>
      <StyledSelect value={selectedValue} onChange={updateSelectedOption}>
        {options.map(option => (
          <MenuItem key={option} value={option}>
            {option}
          </MenuItem>
        ))}
      </StyledSelect>
    </FormControl>
  );
};

export default Select;
