import * as React from "react";
import {
  FormControl as MuiFormControl,
  FormControlLabel,
  FormLabel as MuiFormLabel,
  RadioGroup as MuiRadioGroup,
} from "@material-ui/core";
import styled from "styled-components";

import Radio from "./radio";

const FormLabel = styled(MuiFormLabel)`
  ${({ theme }) => `
  && {
    color: ${theme.palette.text.primary};
  }
  font-weight: bold;
  position: relative;
  &.Mui-disabled {
    opacity: 0.75;
  }
  `}
`;

const FormControl = styled(MuiFormControl)`
  margin: 16px 0;
  min-width: fit-content;
`;

interface RadioGroupOption {
  label: string;
  value?: string;
}

export interface RadioGroupProps {
  defaultOption?: number;
  label?: string;
  disabled?: boolean;
  name: string;
  options: RadioGroupOption[];
  onChange: (value: string) => void;
}

const RadioGroup: React.FC<RadioGroupProps> = ({
  defaultOption = 0,
  label,
  disabled,
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
    <FormControl key={name} disabled={disabled}>
      {label && <FormLabel>{label}</FormLabel>}
      <MuiRadioGroup
        aria-label={label}
        name={name}
        defaultValue={options[defaultIdx]?.value || options[defaultIdx].label}
        onChange={updateSelectedOption}
      >
        {options.map((option, idx) => {
          return (
            <FormControlLabel
              key={option.label}
              value={option?.value || option.label}
              control={<Radio selected={idx === selectedIdx} />}
              label={option.label}
            />
          );
        })}
      </MuiRadioGroup>
    </FormControl>
  );
};

export default RadioGroup;
