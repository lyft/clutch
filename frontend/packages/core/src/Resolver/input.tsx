import React from "react";
import type { UseFormMethods } from "react-hook-form";
import type { clutch } from "@clutch-sh/api";
import { useForm } from "react-hook-form";


import { Error } from "../error";
import { TextField } from "../Input/text-field";

import type { ChangeEventTarget, ResolverChangeEvent } from "./hydrator";
import { convertChangeEvent, hydrateField } from "./hydrator";
import { Accordion, AccordionDetails, AccordionDivider, AccordionActions } from "../accordion";
import { Button } from "../button";

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
interface SchemaResolverProps {
  schemas: clutch.resolver.v1.Schema[];
  selectedSchema: number;
  onChange: (e: ResolverChangeEvent) => void;
  validation: UseFormMethods;
}

const SchemaResolver = ({ schema, expanded, onClick }) => {
  const schemaValidation = useForm({
    mode: "onSubmit",
    reValidateMode: "onSubmit",
    shouldFocusError: false,
  });

  const onChange = (e) => { console.log("schema event", JSON.stringify(e)) }; // TODO

  return (
    <form>
      <Accordion
        title={`Search by ${schema.metadata.displayName}`}
        expanded={expanded}
        onClick={onClick}
      >
        <AccordionDetails>
          {schema.error ? (
            <Error message={`Schema Error: ${schema.error.message}`} />
          ) :
            schema.fields.map(field => hydrateField(field, onChange, schemaValidation))
          }
        </AccordionDetails>
        <AccordionDivider />
        <AccordionActions>
          <Button text="Submit" />
        </AccordionActions>
      </Accordion>
    </form>
  );
};


export { SchemaResolver, QueryResolver };
