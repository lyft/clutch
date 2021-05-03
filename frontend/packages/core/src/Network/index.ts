import type { AxiosError, AxiosResponse } from "axios";
import axios from "axios";

import type { ClutchError } from "./errors";
import { grpcResponseToError } from "./errors";

/**
 * HTTP response status.
 *
 * Responses are grouped in five classes:
 *  - Informational responses (100–199)
 *  - Successful responses (200–299)
 *  - Redirects (300–399)
 *  - Client errors (400–499)
 *  - Server errors (500–599)
 */
export interface HttpStatus {
  /** The status code. */
  code: number;
  /** The status message. */
  text: string;
}

const successInterceptor = (response: AxiosResponse) => {
  return response;
};

const errorInterceptor = (error: AxiosError): Promise<ClutchError> => {
  const response = error?.response;
  if (response === undefined) {
    const clientError = {
      status: {
        code: 500,
        text: "Client Error",
      },
      message: error.message,
    } as ClutchError;
    return Promise.reject(clientError);
  }

  // This section handles authentication redirects.
  if (response?.status === 401) {
    // TODO: turn this in to silent refresh once refresh tokens are supported.
    const redirectUrl = window.location.pathname + window.location.search;
    window.location.href = `/v1/authn/login?redirect_url=${encodeURIComponent(redirectUrl)}`;
  }

  // we are guaranteed to have a response object on the error from this point on
  // since we have already accounted for axios errors.
  const responseData = error?.response?.data;
  // if the response data has a code on it we know it's a gRPC response.
  let err;
  if (responseData?.code !== undefined) {
    err = grpcResponseToError(error);
  } else {
    const message =
      typeof error.response?.data === "string"
        ? error.response.data
        : error?.message || error.response.statusText;
    err = {
      status: {
        code: error.response.status,
        text: error.response.statusText,
      } as HttpStatus,
      message,
      data: responseData,
    } as ClutchError;
  }
  return Promise.reject(err);
};

const createClient = () => {
  const axiosClient = axios.create({
    // n.b. the client will treat any response code >= 400 as an error and apply the error interceptor.
    validateStatus: status => {
      return status < 400;
    },
  });
  axiosClient.interceptors.response.use(successInterceptor, errorInterceptor);

  return axiosClient;
};

const client = createClient();

export { client, errorInterceptor, successInterceptor };
