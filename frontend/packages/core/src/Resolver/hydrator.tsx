import React from "react";
import type { clutch } from "@clutch-sh/api";
import _ from "lodash";

import Select from "../Input/select";
import TextField from "../Input/text-field";

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
  const { name, metadata } = field;
  const errorMsg = validation?.errors?.[name]?.message || validation?.errors?.[name]?.type || "";

  const handleChanges = (event: React.ChangeEvent<ChangeEventTarget> | React.KeyboardEvent) => {
    onChange(convertChangeEvent(event));
  };

  return (
    <TextField
      key={metadata.displayName || name}
      placeholder={metadata.stringField.placeholder}
      defaultValue={metadata.stringField.defaultValue || null}
      required={metadata.required || false}
      name={name}
      label={metadata.displayName || name}
      onChange={handleChanges}
      onKeyDown={handleChanges}
      onFocus={handleChanges}
      inputRef={validation.register({ required: metadata.required || false })}
      helperText={errorMsg}
      error={!!errorMsg}
    />
  );
};

const OptionField = (
  field: clutch.resolver.v1.IField,
  onChange: (e: ResolverChangeEvent) => void
): React.ReactElement => {
  const { name, metadata } = field;
  const sortedOptions = _.sortBy(metadata.optionField.options, o => o.displayName);
  React.useEffect(() => {
    onChange({
      target: {
        name,
        value: sortedOptions?.[0]?.stringValue,
      },
      initialLoad: true,
    });
  }, []);

  const options = sortedOptions.map(option => {
    return { label: option.displayName, value: option.stringValue };
  });
  const updateSelectedOption = (value: string) => {
    onChange({
      target: {
        name,
        value,
      },
    });
  };

  return (
    <Select
      key={metadata.displayName}
      label={metadata.displayName}
      onChange={updateSelectedOption}
      name={name}
      options={options}
    />
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

export { convertChangeEvent, hydrateField };
