import React from "react";
import type { UseFormMethods } from "react-hook-form";
import { useForm } from "react-hook-form";
import type { clutch } from "@clutch-sh/api";

import {
  Accordion,
  AccordionActions,
  AccordionDetails,
  AccordionDivider,
  AccordionProps,
} from "../accordion";
import { Button } from "../button";
import { Error } from "../error";
import { TextField } from "../Input/text-field";

import type { ChangeEventTarget, ResolverChangeEvent } from "./hydrator";
import { convertChangeEvent, hydrateField } from "./hydrator";

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

// TODO: update and use
interface SchemaResolverProps extends Pick<AccordionProps, "expanded" | "onClick"> {
  schema: clutch.resolver.v1.Schema;
  submitHandler: any;
}

const SchemaResolver = ({ schema, expanded, onClick, submitHandler }: SchemaResolverProps) => {
  const [data, setData] = React.useState({"@type" : schema.typeUrl});

  const schemaValidation = useForm({
    mode: "onSubmit",
    reValidateMode: "onSubmit",
    shouldFocusError: false,
  });

  const onChange = e => {
    setData({ ...data, [e.target.name]: e.target.value });
  };

  return (
    <form noValidate onSubmit={schemaValidation.handleSubmit(() => submitHandler(data))}>
      <Accordion
        title={`Search by ${schema.metadata.displayName}`}
        expanded={expanded}
        onClick={onClick}
      >
        <AccordionDetails>
          {schema.error ? (
            <Error message={`Schema Error: ${schema.error.message}`} />
          ) : (
            schema.fields.map(field => hydrateField(field, onChange, schemaValidation))
          )}
        </AccordionDetails>
        <AccordionDivider />
        <AccordionActions>
          <Button text="Submit" type="submit" />
        </AccordionActions>
      </Accordion>
    </form>
  );
};

export { SchemaResolver, QueryResolver };
