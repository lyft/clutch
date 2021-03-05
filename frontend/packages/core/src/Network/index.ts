import type { AxiosError } from "axios";
import axios from "axios";

import type { ClutchError } from "./errors";
import grpcResponseToError from "./errors";

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

const errorInterceptor = (error: AxiosError): Promise<ClutchError> => {
  if (error.isAxiosError) {
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
  const responseData = error?.response.data;
  // if the response data has a code on it we know it's a gRPC response.
  let err;
  if (responseData?.code !== undefined) {
    err = grpcResponseToError(error);
  } else {
    err = {
      status: {
        code: error.response.status,
        text: error.response.statusText,
      } as HttpStatus,
      message: error?.message || error.response.statusText,
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
  axiosClient.interceptors.response.use(resp => resp, errorInterceptor);

  return axiosClient;
};

const client = createClient();

export { client, errorInterceptor };
