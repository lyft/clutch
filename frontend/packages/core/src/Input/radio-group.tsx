import React from "react";
import {
  FormControl as MuiFormControl,
  FormControlLabel,
  FormLabel as MuiFormLabel,
  Radio as MuiRadio,
  RadioGroup as MuiRadioGroup,
} from "@material-ui/core";
import styled from "styled-components";

const FormLabel = styled(MuiFormLabel)`
  ${({ theme }) => `
  && {
    color: ${theme.palette.text.primary};
  }
  font-weigh: bold;
  position: relative;
  `}
`;

const StyledRadioGroup = styled(MuiRadioGroup)`
  ${({ ...props }) => `
  display: flex;
  max-width: ${props["data-max-width"] || "500px"};
  `}
`;

const FormControl = styled(MuiFormControl)`
  ${({ ...props }) => `
  display: flex;
  min-width: fit-content;
  width: ${props["data-max-width"] || "500px"};
  `}
`;

const Radio = styled(MuiRadio)`
  ${({ theme }) => `
  &.Mui-checked {
    color: ${theme.palette.accent.main};
  }
  `}
`;

interface RadioGroupOption {
  label: string;
  value?: string;
}

interface RadioGroupProps {
  defaultOption?: number;
  label?: string;
  maxWidth?: string;
  name: string;
  options: RadioGroupOption[];
  onChange: (value: string) => void;
}

const RadioGroup: React.FC<RadioGroupProps> = ({
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
    const optionValues = options.map(o => o?.value || o.label);
    setSelectedIdx(optionValues.indexOf(value));
    if (onChange !== undefined) {
      onChange(value);
    }
  };

  React.useEffect(() => {
    if (onChange !== undefined) {
      onChange(options[selectedIdx]?.value || options[selectedIdx].label);
    }
  }, []);

  return (
    <FormControl key={name} data-max-width={maxWidth}>
      {label && <FormLabel>{label}</FormLabel>}
      <StyledRadioGroup
        aria-label={label}
        name={name}
        defaultValue={options[defaultIdx]?.value || options[defaultIdx].label}
        onChange={updateSelectedOption}
      >
        {options.map(option => {
          return (
            <FormControlLabel
              key={option.label}
              value={option?.value || option.label}
              control={<Radio />}
              label={option.label}
            />
          );
        })}
      </StyledRadioGroup>
    </FormControl>
  );
};

export { RadioGroup, RadioGroupProps };
