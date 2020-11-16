import React from "react";
import { useForm } from "react-hook-form";
import { DevTool } from "@hookform/devtools";
import _ from "lodash";
import styled from "styled-components";

import { Button } from "../button";
import { useWizardContext } from "../Contexts";
import { CompressedError, Error } from "../error";
import Loadable from "../loading";

import { fetchResourceSchemas, resolveResource } from "./fetch";
import type { ResolverChangeEvent } from "./hydrator";
import { QueryResolver, SchemaResolver } from "./input";
import type { DispatchAction } from "./state";
import { ResolverAction, useResolverState } from "./state";

const Spacer = styled.div`
  margin: 10px;
`;

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

  const queryValidation = useForm({
    mode: "onSubmit",
    reValidateMode: "onSubmit",
    shouldFocusError: false,
  });
  const schemaValidation = useForm({
    mode: "onSubmit",
    reValidateMode: "onSubmit",
    shouldFocusError: false,
  });
  const [validation, setValidation] = React.useState(() => queryValidation);

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

  const updateResolverData = (e: ResolverChangeEvent) => {
    validation.clearErrors();
    if (e.target.name !== "query") {
      setValidation(() => schemaValidation);
    } else {
      setValidation(() => queryValidation);
    }
    dispatch({
      type: ResolverAction.UPDATE_QUERY_DATA,
      data: { [e.target.name.toLowerCase()]: e.target.value },
    });
    if (e.initialLoad) {
      setValidation(() => queryValidation);
    }
  };

  const setSelectedSchema = (e: React.ChangeEvent<{ name?: string; value: unknown }>) => {
    dispatch({ type: ResolverAction.SET_SELECTED_SCHEMA, schema: e.target.value });
  };

  return (
    <Loadable isLoading={state.schemasLoading}>
      {state.schemaFetchError !== "" ? (
        <Error message={state.schemaFetchError} onRetry={() => loadSchemas(type, dispatch)} />
      ) : (
        <Loadable variant="overlay" isLoading={state.resolverLoading}>
          {process.env.REACT_APP_DEBUG_FORMS === "true" && <DevTool control={validation.control} />}
          {(variant === "dual" || variant === "query") && (
            <Form onSubmit={validation.handleSubmit(submitHandler)} noValidate>
              <QueryResolver
                schemas={state.searchableSchemas}
                onChange={updateResolverData}
                validation={queryValidation}
              />
            </Form>
          )}
          {variant === "dual" && (
            <>
              <Spacer />- OR -
            </>
          )}
          {(variant === "dual" || variant === "schema") && (
            <Form onSubmit={validation.handleSubmit(submitHandler)} noValidate>
              <SchemaResolver
                schemas={state.allSchemas}
                selectedSchema={state.selectedSchema}
                onSelect={setSelectedSchema}
                onChange={updateResolverData}
                validation={schemaValidation}
              />
            </Form>
          )}
          <Button text="Continue" onClick={validation.handleSubmit(submitHandler)} />
          <CompressedError title="Error" message={state.resolverFetchError} />
        </Loadable>
      )}
    </Loadable>
  );
};

export default Resolver;
