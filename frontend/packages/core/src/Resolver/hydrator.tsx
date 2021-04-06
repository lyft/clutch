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
  const errorMsg =
    validation?.errors?.[field.name]?.message || validation?.errors?.[field.name]?.type || "";

  const handleChanges = (event: React.ChangeEvent<ChangeEventTarget> | React.KeyboardEvent) => {
    onChange(convertChangeEvent(event));
  };

  return (
    <TextField
      key={field.metadata.displayName || field.name}
      placeholder={field.metadata.stringField.placeholder}
      defaultValue={field.metadata.stringField.defaultValue || null}
      required={field.metadata.required || false}
      name={field.name}
      label={field.metadata.displayName || field.name}
      onChange={handleChanges}
      onKeyDown={handleChanges}
      onFocus={handleChanges}
      helperText={errorMsg}
      error={!!errorMsg}
      {...validation.register(field.name, { required: field.metadata.required || false })}
    />
  );
};

const OptionField = (
  field: clutch.resolver.v1.IField,
  onChange: (e: ResolverChangeEvent) => void
): React.ReactElement => {
  React.useEffect(() => {
    onChange({
      target: {
        name: field.name,
        value: field.metadata.optionField.options?.[0]?.stringValue,
      },
      initialLoad: true,
    });
  }, []);

  const options = _.sortBy(field.metadata.optionField.options, o => o.displayName).map(option => {
    return { label: option.displayName, value: option.stringValue };
  });
  const updateSelectedOption = (value: string) => {
    onChange({
      target: {
        name: field.name,
        value,
      },
    });
  };

  return (
    <Select
      key={field.metadata.displayName}
      label={field.metadata.displayName}
      onChange={updateSelectedOption}
      name={field.name}
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
