import React from "react";
import styled from "@emotion/styled";
import _ from "lodash";

import { AccordionGroup } from "../accordion";
import { useWizardContext } from "../Contexts";
import { CompressedError, Error } from "../Feedback";
import { HorizontalRule } from "../horizontal-rule";
import Loadable from "../loading";

import { fetchResourceSchemas, resolveResource } from "./fetch";
import { QueryResolver, SchemaResolver } from "./input";
import type { DispatchAction } from "./state";
import { ResolverAction, useResolverState } from "./state";

const SchemaLabel = styled.div({
  alignSelf: "flex-start",
  fontSize: "20px",
  fontWeight: 700,
  marginBottom: "8px",
});

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
  /**
   *  API module to resolve lookups against.
   * */
  apiPackage?: object;
}

const Resolver: React.FC<ResolverProps> = ({
  type,
  searchLimit,
  onResolve,
  variant = "dual",
  apiPackage,
}) => {
  const [state, dispatch] = useResolverState();
  const { displayWarnings } = useWizardContext();

  React.useEffect(() => loadSchemas(type, dispatch), []);

  const submitHandler = data => {
    // Move to loading state.
    dispatch({ type: ResolverAction.RESOLVING });

    // Copy incoming data, trimming whitespace from any string values (usually artifact of cut and paste into tool).
    const inputData = _.mapValues(data, v => (_.isString(v) && _.trim(v)) || v);

    // Resolve!
    resolveResource(
      type,
      searchLimit,
      inputData,
      (results, failures) => {
        onResolve({ results, input: inputData });
        if (!_.isEmpty(failures)) {
          displayWarnings(failures);
        }
        dispatch({ type: ResolverAction.RESOLVE_SUCCESS });
      },
      err => dispatch({ type: ResolverAction.RESOLVE_ERROR, error: err }),
      apiPackage
    );
  };

  return (
    <Loadable isLoading={state.schemasLoading}>
      {state.schemaFetchError !== "" ? (
        <Error message={state.schemaFetchError} onRetry={() => loadSchemas(type, dispatch)} />
      ) : (
        <Loadable variant="overlay" isLoading={state.resolverLoading}>
          <CompressedError title="Error" message={state.resolverFetchError} />
          {(variant === "dual" || variant === "query") && (
            <>
              <SchemaLabel>Search</SchemaLabel>
              <QueryResolver schemas={state.searchableSchemas} submitHandler={submitHandler} />
            </>
          )}
          {variant === "dual" && <HorizontalRule>OR</HorizontalRule>}
          <SchemaLabel>Advanced Search</SchemaLabel>
          <AccordionGroup defaultExpandedIdx={0}>
            {state.allSchemas.map(schema => (
              <SchemaResolver key={schema.typeUrl} schema={schema} submitHandler={submitHandler} />
            ))}
          </AccordionGroup>
        </Loadable>
      )}
    </Loadable>
  );
};

export default Resolver;
