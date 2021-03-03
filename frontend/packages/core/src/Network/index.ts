import type { AxiosError } from "axios";
import axios from "axios";

import type { ClientError, ClutchError, NetworkError } from "./errors";
import clutchError from "./errors";

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

const errorInterceptor = (error: AxiosError): Promise<ClutchError | NetworkError | ClientError> => {
  if (error.isAxiosError) {
    const clientError = { message: `Client Error: ${error.message}` } as ClientError;
    return Promise.reject(clientError);
  }

  // n.b. we are guaranteed to have a response object on the error from this point on
  // since we have already accounted for axios errors.

  // TODO: check for authentication errors and handle
  const responseData = error?.response.data;
  // n.b. if the response data has a code on it we know it's a gRPC response.
  let err;
  if (responseData?.code !== undefined) {
    err = clutchError(error);
  } else {
    err = {
      status: {
        code: error.response.status,
        text: error.response.statusText,
      } as HttpStatus,
      message: error?.message || error.response.statusText,
      data: responseData,
    } as NetworkError;
  }

  return Promise.reject(err);
};

const createClient = () => {
  const axiosClient = axios.create({
    validateStatus: status => {
      return status < 400;
    },
  });
  axiosClient.interceptors.response.use(resp => resp, errorInterceptor);

  return axiosClient;
};

const client = createClient();

export { client, errorInterceptor };
