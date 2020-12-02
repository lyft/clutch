import React from "react";
import type { UseFormMethods } from "react-hook-form";
import type { clutch } from "@clutch-sh/api";
import { MenuItem, Select } from "@material-ui/core";

import { Error } from "../error";
import { TextField } from "../Input/text-field";

import type { ChangeEventTarget, ResolverChangeEvent } from "./hydrator";
import { convertChangeEvent, FormControl, hydrateField, InputLabel } from "./hydrator";

interface QueryResolverProps {
  schemas: clutch.resolver.v1.Schema[];
  onChange: (e: ResolverChangeEvent) => void;
  validation: UseFormMethods;
}

const QueryResolver: React.FC<QueryResolverProps> = ({ schemas, onChange, validation }) => {
  let typeLabel = schemas.map(schema => schema?.metadata.displayName).join();
  typeLabel = `Search by ${typeLabel}`;

  const handleChanges = (event: React.ChangeEvent<ChangeEventTarget> | React.KeyboardEvent) => {
    onChange(convertChangeEvent(event));
  };

  const error = validation.errors?.query;
  return (
    <TextField
      label={typeLabel || "Please select a resolver"}
      name="query"
      required
      onChange={handleChanges}
      onKeyDown={handleChanges}
      onFocus={handleChanges}
      inputRef={validation.register({ required: true })}
      error={!!error}
      helperText={error?.message || error?.type || ""}
    />
  );
};

interface SchemaResolverProps {
  schemas: clutch.resolver.v1.Schema[];
  selectedSchema: number;
  onSelect: (e: React.ChangeEvent<{ name?: string; value: unknown }>) => void;
  onChange: (e: ResolverChangeEvent) => void;
  validation: UseFormMethods;
}

const SchemaResolver: React.FC<SchemaResolverProps> = ({
  schemas,
  selectedSchema,
  onSelect,
  onChange,
  validation,
}) => (
  <>
    <FormControl>
      <InputLabel>Resolver</InputLabel>
      <Select value={schemas?.[selectedSchema]?.typeUrl || ""} onChange={onSelect}>
        {schemas.map(schema => (
          <MenuItem key={schema.metadata.displayName} value={schema.typeUrl}>
            {schema.metadata.displayName}
          </MenuItem>
        ))}
      </Select>
    </FormControl>
    {schemas[selectedSchema]?.error ? (
      <Error message={`Schema Error: ${schemas[selectedSchema].error.message}`} />
    ) : (
      schemas[selectedSchema]?.fields.map(field => hydrateField(field, onChange, validation))
    )}
  </>
);

export { SchemaResolver, QueryResolver };
