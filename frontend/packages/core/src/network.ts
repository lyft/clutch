import type { AxiosError, AxiosResponse } from "axios";
import axios from "axios";

interface HttpStatus {
  code: number;
  text: string;
}

// Add status codes here as needed when looking at code outside of gRPC context.
const HTTP_STATUS_MAPPINGS = {
  UNAUTHORIZED: 401,
};

// Pulled from Google's mapping:
// https://github.com/grpc-ecosystem/grpc-gateway/blob/e1a127c6f7cf2006c77a16dbc0face2a43f2094a/third_party/googleapis/google/rpc/code.proto
const GRPC_CODE_MAPPINGS = {
  0: { code: 200, text: "OK" },
  1: { code: 499, text: "Cancelled" },
  2: { code: 500, text: "Internal Server Error" },
  3: { code: 400, text: "Invalid Argument" },
  4: { code: 504, text: "Gateway Timeout" },
  5: { code: 404, text: "Not Found" },
  6: { code: 409, text: "Already Exists" },
  7: { code: 403, text: "Permission Denied" },
  8: { code: 429, text: "Resource Exhausted" },
  9: { code: 400, text: "Failed Precondition" },
  10: { code: 409, text: "Aborted" },
  11: { code: 400, text: "Out-of-Range" },
  12: { code: 501, text: "Not Implemented" },
  13: { code: 500, text: "Internal Server Error" },
  14: { code: 503, text: "Service Unavailable" },
  15: { code: 500, text: "Internal Server Error" },
  16: { code: 401, text: "Unauthenticated" },
};

const DEFAULT_CODE = { code: 500, text: "Unknown Code" };

const grpcCodeToHttpStatus = (code: number): HttpStatus => {
  const convertedCode = GRPC_CODE_MAPPINGS?.[code];
  if (convertedCode === undefined) {
    return DEFAULT_CODE;
  }
  return convertedCode;
};

interface ClientErrorMessage {
  summary: string;
  details: string;
}

/**
 * Parse an error message for the summary and details.
 *
 * The summary is captured on the first line of the message
 * and the details are everything following.
 *
 * @param {*} error
 */
const parseErrorMessage = (error: string): ClientErrorMessage => {
  let detailBreakpoint = error.indexOf("\n");
  detailBreakpoint = detailBreakpoint === -1 ? error.length : detailBreakpoint;
  const summary = error.substring(0, detailBreakpoint);
  const details = error.substring(detailBreakpoint + 1);
  return { summary, details };
};

interface Response extends AxiosResponse {
  statusCode: number;
  statusText: string;
  displayText: string;
  details: string;
}

/**
 * Represents a Client Error with custom props.
 *
 * @returns {ClientError} Returns a client error.
 */
class ClientError extends Error {
  response: Response;

  constructor(resp: AxiosResponse, ...params: any[]) {
    super(...params);
    const { summary, details } = parseErrorMessage(resp.data.message);
    this.response = {
      ...resp,
      statusCode: grpcCodeToHttpStatus(resp.data.code).code,
      statusText: grpcCodeToHttpStatus(resp.data.code).text,
      displayText: `${grpcCodeToHttpStatus(resp.data.code).text}: ${summary}`,
      details,
    };
  }
}

const successInterceptor = (response: AxiosResponse<any>) => {
  // n.b. this middleware handles authentication redirects
  // to prevent CORS issues from redirecting on the server.
  if (response?.data?.authUrl) {
    window.location = response.data.authUrl;
    response.data.code = 401;
    response.data.message = "Authentication Expired";
    throw new ClientError(response);
  }

  return response;
};

const errorInterceptor = (error: AxiosError) => {
  const { response } = error;

  if (response?.status === HTTP_STATUS_MAPPINGS.UNAUTHORIZED) {
    const dest = encodeURIComponent(window.location.pathname);
    window.location.href = `${window.location.origin}/v1/authn/login?redirect_url=${dest}`;
  }

  if (response.data?.code && response.data?.message) {
    throw new ClientError(response);
  }
  // We tried! Return a failed promise.
  return Promise.reject(error);
};

const createClient = () => {
  const instance = axios.create();
  instance.interceptors.response.use(
    response => successInterceptor(response),
    error => errorInterceptor(error)
  );

  return instance;
};

const client = createClient();

export { client, ClientError, grpcCodeToHttpStatus, parseErrorMessage };
