import React from "react";
import { useForm } from "react-hook-form";
import { DevTool } from "@hookform/devtools";
import _ from "lodash";
import styled from "styled-components";

import { Button } from "../button";
import { useWizardContext } from "../Contexts";
import { CompressedError, Error } from "../error";
import { HorizontalRule } from "../horizontal-rule";
import Loadable from "../loading";

import { fetchResourceSchemas, resolveResource } from "./fetch";
import { hydrateField, ResolverChangeEvent } from "./hydrator";
import { QueryResolver, SchemaResolver } from "./input";
import type { DispatchAction } from "./state";
import { ResolverAction, useResolverState } from "./state";
import { Accordion, AccordionActions, AccordionDetails, AccordionGroup } from "../accordion";
import TextField from "../Input/text-field";
import Select from "../Input/select";

const Form = styled.form`
  align-items: center;
  display: flex;
  flex-direction: column;
  width: 100%;
`;

const loadSchemas = (type: string, dispatch: React.Dispatch<DispatchAction>) => {
  fetchResourceSchemas(type)
    .then(schemas => {
      if (schemas.length === 0) {
        dispatch({
          type: ResolverAction.SCHEMAS_ERROR,
          error: `No schemas found for type '${type}'`,
        });
      } else {
        dispatch({ type: ResolverAction.SCHEMAS_SUCCCESS, allSchemas: schemas });
      }
    })
    .catch(err => {
      dispatch({ type: ResolverAction.SCHEMAS_ERROR, error: err.message });
    });
};

interface ResolverProps {
  type: string;
  searchLimit: number;
  onResolve: (data: { results: object[]; input: object }) => void;
  variant?: "dual" | "query" | "schema";
}

const Resolver: React.FC<ResolverProps> = ({ type, searchLimit, onResolve, variant = "dual" }) => {
  const [state, dispatch] = useResolverState();
  const { displayWarnings } = useWizardContext();

  React.useEffect(() => loadSchemas(type, dispatch), []);

  const submitHandler = () => {
    // Move to loading state.
    dispatch({ type: ResolverAction.RESOLVING });

    // Copy incoming data, trimming whitespace from any string values (usually artifact of cut and paste into tool).
    const data = _.mapValues(state.queryData, v => (_.isString(v) && _.trim(v)) || v);

    // Set desired type.
    data["@type"] = state.allSchemas[state.selectedSchema]?.typeUrl;

    // Resolve!
    resolveResource(
      type,
      searchLimit,
      data,
      (results, failures) => {
        onResolve({ results, input: data });
        if (!_.isEmpty(failures)) {
          displayWarnings(failures);
        }
        dispatch({ type: ResolverAction.RESOLVE_SUCCESS });
      },
      err => dispatch({ type: ResolverAction.RESOLVE_ERROR, error: err })
    );
  };



  const queryValidation = useForm({
    mode: "onSubmit",
    reValidateMode: "onSubmit",
    shouldFocusError: false,
  });

  return (
    <Loadable isLoading={state.schemasLoading}>
      {state.schemaFetchError !== "" ? (
        <Error message={state.schemaFetchError} onRetry={() => loadSchemas(type, dispatch)} />
      ) : (
          <Loadable variant="overlay" isLoading={state.resolverLoading}>
            {/* {process.env.REACT_APP_DEBUG_FORMS === "true" && <DevTool control={validation.control} />} */}
            {(variant === "dual" || variant === "query") && (
              <Form onSubmit={queryValidation.handleSubmit(submitHandler)} noValidate>
                <QueryResolver
                  schemas={state.searchableSchemas}
                  onChange={() => {}}
                  validation={queryValidation}
                />
                <Button text="Search" />
              </Form>
            )}
            {variant === "dual" && <HorizontalRule>OR</HorizontalRule>}
          Advanced Search
            <AccordionGroup defaultExpandedIdx={0}>
              {state.allSchemas.map((schema, idx) => <SchemaResolver2 key={schema.typeUrl} schema={schema} />)}
            </AccordionGroup>

          </Loadable>
        )}
    </Loadable>
  );
};

const SchemaResolver2 = ({ schema, expanded, onClick }) => {
  console.log(schema);

  return <Accordion title={`Search by ${schema.metadata.displayName}`} expanded={expanded} onClick={onClick}>
    <AccordionDetails>
      {schema.error ? (<Error message={`Schema Error: ${schema.error.message}`} />) :
        (schema?.fields.map(field => {
          if (field.metadata.type === "stringField") {
            return <TextField key={field.metadata.name} label={field.metadata.displayName} />
          } else if (field.metadata.type === "optionField") {
            const options = field.metadata.optionField.options.map(option => {
              return { label: option.displayName, value: option.stringValue };
            });

            return <Select key={field.metadata.name} label={field.metadata.displayName} options={options} onChange={() => { }} />
          }
        }))
      }
    </AccordionDetails>
    <AccordionActions>
      <Button text="Submit" />
    </AccordionActions>
  </Accordion>

};

export default Resolver;
