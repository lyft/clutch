import type { clutch as IClutch } from "@clutch-sh/api";
import * as $pbclutch from "@clutch-sh/api";
import type { ClutchError } from "@clutch-sh/core";
import _ from "lodash";

import { client } from "../Network";

const fetchResourceSchemas = async (type: string): Promise<IClutch.resolver.v1.Schema[]> => {
  const response = await client.post("/v1/resolver/getObjectSchemas", {
    type_url: `type.googleapis.com/${type}`,
  });
  return response.data.schemas.map((schema: object) =>
    $pbclutch.clutch.resolver.v1.Schema.fromObject(schema)
  );
};

export interface ResolutionResults {
  results: object[];
  failures: {
    message: string;
  }[];
}

const resolveQuery = async (
  type: string,
  limit: number,
  fields: {
    [key: string]: any;
  }
): Promise<ResolutionResults> => {
  const response = await client.post("/v1/resolver/search", {
    want: `type.googleapis.com/${type}`,
    query: fields.query,
    limit,
  });
  return { results: response.data.results, failures: response.data.partialFailures };
};

const resolveFields = async (
  type: string,
  limit: number,
  fields: object
): Promise<ResolutionResults> => {
  const response = await client.post("/v1/resolver/resolve", {
    want: `type.googleapis.com/${type}`,
    have: fields,
    limit,
  });
  return { results: response.data?.results || [], failures: response.data?.partialFailures || [] };
};

const resolveResource = async (
  type: string,
  limit: number,
  fields: {
    [key: string]: any;
  },
  onResolve: (resultObjects: any[], failureMessages: string[]) => void,
  onError: (message: ClutchError) => void,
  apiPackage?: any
) => {
  const resolver = fields?.query !== undefined ? resolveQuery : resolveFields;
  return resolver(type, limit, fields)
    .then(({ results, failures }) => {
      // n.b. default to using the open source @clutch-sh/api package to resolve the
      // resource against unless a custom package has been specified by the workflow.
      let pbClutch = _.get($pbclutch, type);
      if (apiPackage) {
        pbClutch = _.get(apiPackage, type);
      }
      const resultObjects = results.map(result => pbClutch.fromObject(result));
      const partialFailures = failures.map(failure => failure.message);
      if (_.some(resultObjects) !== undefined) {
        onResolve(resultObjects, partialFailures);
      }
    })
    .catch(err => onError(err));
};

export { fetchResourceSchemas, resolveResource };
