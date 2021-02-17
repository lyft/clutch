import React from "react";
import { useForm } from "react-hook-form";
import type { clutch } from "@clutch-sh/api";
import styled from "@emotion/styled";
import SearchIcon from "@material-ui/icons/Search";

import {
  Accordion,
  AccordionActions,
  AccordionDetails,
  AccordionDivider,
  AccordionProps,
} from "../accordion";
import { Button } from "../button";
import { Error } from "../Feedback";
import { TextField } from "../Input/text-field";
import { client } from "../network";

import type { ChangeEventTarget } from "./hydrator";
import { convertChangeEvent, hydrateField } from "./hydrator";

const Form = styled.form({});

interface QueryResolverProps {
  // The inputType is the orignal resource type requested
  // eg: clutch.aws.ec2.v1.AutoscalingGroup
  inputType: string;
  schemas: clutch.resolver.v1.Schema[];
  submitHandler: any;
}

const autoComplete = async (type: string, search: string): Promise<any> => {
  const response = await client.post("/v1/resolver/autocomplete", {
    want: `type.googleapis.com/${type}`,
    search,
  });

  return { results: response?.data?.results || [] };
};

const QueryResolver: React.FC<QueryResolverProps> = ({ inputType, schemas, submitHandler }) => {
  const validation = useForm({
    mode: "onSubmit",
    reValidateMode: "onSubmit",
    shouldFocusError: false,
  });

  const [queryData, setQueryData] = React.useState("");

  let typeLabel = schemas.map(schema => schema?.metadata.displayName).join();
  typeLabel = `Search by ${typeLabel}`;

  const handleChanges = (event: React.ChangeEvent<ChangeEventTarget> | React.KeyboardEvent) => {
    setQueryData(convertChangeEvent(event).target.value);
  };

  // If there is at least 1 schema that has the ability to autocomplete we will enable it.
  const isAutoCompleteEnabled =
    schemas.filter(schema => schema?.metadata?.search?.autocompleteEnabled === true).length >= 1;

  const error = validation.errors?.query;
  return (
    <Form onSubmit={validation.handleSubmit(() => submitHandler({ query: queryData }))} noValidate>
      <TextField
        label={typeLabel}
        name="query"
        required
        onChange={handleChanges}
        onKeyDown={handleChanges}
        onFocus={handleChanges}
        inputRef={validation.register({ required: true })}
        error={!!error}
        helperText={error?.message || error?.type || ""}
        endAdornment={<SearchIcon />}
        enableAutocomplete={isAutoCompleteEnabled}
        autocompleteCallback={v => autoComplete(inputType, v)}
      />
    </Form>
  );
};

// TODO: update and use
interface SchemaResolverProps extends Pick<AccordionProps, "expanded" | "onClick"> {
  schema: clutch.resolver.v1.Schema;
  submitHandler: any;
}

const SchemaDetails = styled(AccordionDetails)({
  "> *": {
    flex: "1 50%",
  },
});

const SchemaResolver = ({ schema, expanded, onClick, submitHandler }: SchemaResolverProps) => {
  const [data, setData] = React.useState({ "@type": schema.typeUrl });

  const schemaValidation = useForm({
    mode: "onSubmit",
    reValidateMode: "onSubmit",
    shouldFocusError: false,
  });

  const onChange = e => {
    setData(existing => {
      return { ...existing, [e.target.name]: e.target.value };
    });
  };

  return (
    <Form noValidate onSubmit={schemaValidation.handleSubmit(() => submitHandler(data))}>
      <Accordion
        title={`Search by ${schema.metadata.displayName}`}
        expanded={expanded}
        onClick={onClick}
      >
        <SchemaDetails>
          {schema.error ? (
            <Error message={`Schema Error: ${schema.error.message}`} />
          ) : (
            schema.fields.map(field => hydrateField(field, onChange, schemaValidation))
          )}
        </SchemaDetails>
        <AccordionDivider />
        <AccordionActions>
          <Button text="Submit" type="submit" />
        </AccordionActions>
      </Accordion>
    </Form>
  );
};

export { SchemaResolver, QueryResolver };
