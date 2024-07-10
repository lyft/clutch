import React from "react";
import { useForm } from "react-hook-form";
import type { clutch } from "@clutch-sh/api";
import styled from "@emotion/styled";
import { DevTool } from "@hookform/devtools";
import SearchIcon from "@mui/icons-material/Search";
import { Stack } from "@mui/material";

import {
  Accordion,
  AccordionActions,
  AccordionDetails,
  AccordionDivider,
  AccordionProps,
} from "../accordion";
import { Button } from "../button";
import { Alert } from "../Feedback";
import { TextField } from "../Input/text-field";
import { useSearchParams } from "../navigation";
import { client } from "../Network";

import type { ChangeEventTarget } from "./hydrator";
import { convertChangeEvent, hydrateField } from "./hydrator";

const Form = styled.form({});

interface QueryResolverProps {
  /**
   * The inputType is the orignal resource type requested
   * eg: clutch.aws.ec2.v1.AutoscalingGroup
   */
  inputType: string;
  schemas: clutch.resolver.v1.Schema[];
  submitHandler: any;
  enableAutocomplete?: boolean;
}

const autoComplete = async (type: string, search: string): Promise<any> => {
  // Check the length of the search query as the user might empty out the search
  // which will still trigger the on change handler
  if (search.length === 0) {
    return { results: [] };
  }

  const response = await client.post("/v1/resolver/autocomplete", {
    want: `type.googleapis.com/${type}`,
    search,
    caseSensitive: false,
  });

  return { results: response?.data?.results || [] };
};

const QueryResolver: React.FC<QueryResolverProps> = ({
  inputType,
  schemas,
  submitHandler,
  enableAutocomplete = true,
}) => {
  const validation = useForm({
    mode: "onSubmit",
    reValidateMode: "onSubmit",
    shouldFocusError: false,
  });

  const [searchParams] = useSearchParams();
  const [queryData, setQueryData] = React.useState(searchParams.get("q") || "");

  let typeLabel = schemas.map(schema => schema?.metadata.displayName).join();
  typeLabel = `Search by ${typeLabel}`;

  const handleChanges = (event: React.ChangeEvent<ChangeEventTarget> | React.KeyboardEvent) => {
    setQueryData(convertChangeEvent(event).target.value);
  };

  // If there is at least 1 schema that has the ability to autocomplete we will enable it.
  // enableAutocomplete's default value is true.  We only use it (set it to false) when we want to override and disable autocomplete at the workflow level rather than the schema level.
  const isAutoCompleteEnabled =
    schemas.filter(schema => schema?.metadata?.search?.autocompleteEnabled === true).length >= 1 &&
    enableAutocomplete;

  const error = validation.formState.errors?.query;
  return (
    <Form onSubmit={validation.handleSubmit(() => submitHandler({ query: queryData }))} noValidate>
      {import.meta.env.VITE_APP_DEBUG_FORMS && <DevTool control={validation.control} />}
      <TextField
        label={typeLabel}
        name="query"
        required
        defaultValue={queryData}
        onChange={handleChanges}
        onKeyDown={handleChanges}
        onFocus={handleChanges}
        error={!!error}
        helperText={error?.message || error?.type || ""}
        endAdornment={<SearchIcon />}
        autocompleteCallback={isAutoCompleteEnabled ? v => autoComplete(inputType, v) : undefined}
        formRegistration={validation.register}
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
      {/* {import.meta.env.VITE_APP_DEBUG_FORMS && <DevTool control={schemaValidation.control} />} */}
      <Accordion
        title={`Search by ${schema.metadata.displayName}`}
        expanded={expanded}
        onClick={onClick}
      >
        <SchemaDetails>
          {schema.error ? (
            <Alert severity="error">Schema Error: ${schema.error.message}</Alert>
          ) : (
            <Stack spacing={2}>
              {schema.fields.map(field => hydrateField(field, onChange, schemaValidation))}
            </Stack>
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
