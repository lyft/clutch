import type { clutch as IClutch } from "@clutch-sh/api";
import * as $pbclutch from "@clutch-sh/api";
import _ from "lodash";

import { client, parseErrorMessage } from "../network";

const fetchResourceSchemas = async (
  type: string,
  apiProto?: any
): Promise<IClutch.resolver.v1.Schema[]> => {
  const response = await client.post("/v1/resolver/getObjectSchemas", {
    type_url: `type.googleapis.com/${type}`,
  });
  return response.data.schemas.map((schema: object) => {
    if (apiProto) {
      try {
        return apiProto.clutch.resolver.v1.Schema.fromObject(schema);
      } catch {}
    }
    return $pbclutch.clutch.resolver.v1.Schema.fromObject(schema);
  });
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
  onError: (message: string) => void,
  apiProto?: any
) => {
  const resolver = fields?.query !== undefined ? resolveQuery : resolveFields;
  return resolver(type, limit, fields)
    .then(({ results, failures }) => {
      let pbClutch = _.get($pbclutch, type);
      if (apiProto) {
        pbClutch = _.get(apiProto, type);
      }
      const resultObjects = results.map(result => pbClutch.fromObject(result));
      const failureMessages = failures.map(failure => parseErrorMessage(failure.message).summary);
      if (_.some(resultObjects) !== undefined) {
        onResolve(resultObjects, failureMessages);
      }
    })
    .catch(err => {
      if (err?.response === undefined) {
        // Some runtime error we don't know how to handle.
        onError(`Internal Client Error: '${err.message}'. Please contact the workflow developer.`);
        return;
      }

      onError(err.response.displayText || err.response.statusText);
    });
};

export { fetchResourceSchemas, resolveResource };
