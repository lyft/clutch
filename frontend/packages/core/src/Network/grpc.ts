import type { HttpStatus } from "./index";

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
  return convertedCode === undefined ? DEFAULT_CODE : convertedCode;
};

export default grpcCodeToHttpStatus;
