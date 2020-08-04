import React from "react";
import { Controller } from "react-hook-form";
import type { clutch } from "@clutch-sh/api";
import {
  FormControl as MuiFormControl,
  FormHelperText,
  InputLabel as MuiInputLabel,
  MenuItem,
  Select as MuiSelect,
} from "@material-ui/core";
import styled from "styled-components";

import TextField from "../Input/text-field";

const maxWidth = "500px";
const InputLabel = styled(MuiInputLabel)`
  ${({ theme }) => `
    color: ${theme.palette.text.primary};
  `}
`;

const FormControl = styled(MuiFormControl)`
  display: flex;
  width: 100%;
  max-width: ${maxWidth};
`;

const Select = styled(MuiSelect)`
  display: flex;
  width: 100%;
  max-width: ${maxWidth};
`;

export interface ResolverChangeEvent {
  target: {
    name: string;
    value: string;
  };
  initialLoad?: boolean;
}

export interface ChangeEventTarget {
  name: string;
  value: string;
}

const convertChangeEvent = (
  event: React.ChangeEvent<ChangeEventTarget> | React.KeyboardEvent
): ResolverChangeEvent => {
  return {
    target: {
      name: (event.target as ChangeEventTarget).name,
      value: (event.target as ChangeEventTarget).value,
    },
  };
};

const StringField = (
  field: clutch.resolver.v1.IField,
  onChange: (e: ResolverChangeEvent) => void,
  validation: any
): React.ReactElement => {
  const errorMsg =
    validation?.errors?.[field.name]?.message || validation?.errors?.[field.name]?.type || "";

  const handleChanges = (event: React.ChangeEvent<ChangeEventTarget> | React.KeyboardEvent) => {
    onChange(convertChangeEvent(event));
  };

  return (
    <TextField
      color="secondary"
      key={field.metadata.displayName || field.name}
      placeholder={field.metadata.stringField.placeholder}
      defaultValue={field.metadata.stringField.defaultValue || null}
      required={field.metadata.required || false}
      name={field.name}
      label={field.metadata.displayName || field.name}
      onChange={handleChanges}
      onKeyDown={handleChanges}
      onFocus={handleChanges}
      inputRef={validation.register({ required: field.metadata.required || false })}
      helperText={errorMsg}
      error={!!errorMsg}
    />
  );
};

const OptionField = (
  field: clutch.resolver.v1.IField,
  onChange: (e: ResolverChangeEvent) => void,
  validation: any
): React.ReactElement => {
  let options = field.metadata.optionField.options.map(option => {
    return option.displayName;
  });
  const [selectedIdx, setSelectedIdx] = React.useState(0);
  const updateSelectedOption = (event: React.ChangeEvent<ChangeEventTarget>) => {
    setSelectedIdx(options.indexOf(event.target.value));
    onChange(convertChangeEvent(event));
  };

  React.useEffect(() => {
    const fieldName = field.metadata.displayName || field.name;
    onChange({
      target: {
        name: fieldName,
        value: field.metadata.optionField.options?.[selectedIdx]?.stringValue,
      },
      initialLoad: true,
    });
  }, []);

  const requiredWithoutOptions = field.metadata.required && options.length === 0;
  if (requiredWithoutOptions) {
    options = [""];
  }

  const fieldName = (field.metadata.displayName || field.name).toLowerCase();
  return (
    options.length !== 0 && (
      <FormControl
        key={field.metadata.displayName || field.name}
        required={field.metadata.required || false}
        error={validation.errors?.[fieldName] !== undefined || false}
      >
        <InputLabel shrink={options[selectedIdx] !== ""} color="secondary">
          {fieldName}
        </InputLabel>
        <Controller
          control={validation.control}
          name={fieldName}
          defaultValue=""
          rules={{ required: field.metadata.required ? "required" : false }}
          as={
            <Select
              value={options[selectedIdx] || ""}
              onChange={updateSelectedOption}
              inputProps={{
                style: { minWidth: "100px" },
              }}
            >
              {options.map(option => (
                <MenuItem key={option} value={option}>
                  {option}
                </MenuItem>
              ))}
            </Select>
          }
        />
        {validation.errors?.[fieldName] !== undefined && (
          <FormHelperText>{validation.errors[fieldName].message}</FormHelperText>
        )}
      </FormControl>
    )
  );
};

const FIELD_TYPES = {
  stringField: StringField,
  optionField: OptionField,
};

const hydrateField = (
  field: clutch.resolver.v1.IField,
  onChange: (e: ResolverChangeEvent) => void,
  validation: any
) => {
  let component;
  Object.keys(FIELD_TYPES).some(type => {
    if (Object.keys(field.metadata).includes(type)) {
      component = FIELD_TYPES[type];
      return true;
    }
    return false;
  });
  return component(field, onChange, validation);
};

export { convertChangeEvent, FormControl, hydrateField, InputLabel };
