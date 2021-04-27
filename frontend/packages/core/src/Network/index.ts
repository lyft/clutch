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
  // n.b. this middleware handles authentication redirects
  // to prevent CORS issues from redirecting on the server.
  if (response?.data?.authUrl) {
    window.location = response.data.authUrl;
    const clutchError = {
      status: {
        code: 401,
        text: "Authentication Expired",
      },
      message: "Authentication Expired",
    } as ClutchError;
    throw clutchError;
  }

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
