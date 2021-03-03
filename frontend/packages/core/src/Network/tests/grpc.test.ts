import grpcCodeToHttpStatus from "../grpc";

describe("grpcCodeToHttpStatus", () => {
  it("returns default value for undefined code", () => {
    const httpStatus = grpcCodeToHttpStatus(-1);
    expect(httpStatus.code).toEqual(500);
    expect(httpStatus.text).toEqual("Unknown Code");
  });

  it("returns a code and text for valid code", () => {
    const httpStatus = grpcCodeToHttpStatus(0);
    expect(httpStatus.code).toEqual(200);
    expect(httpStatus.text).toEqual("OK");
  });
});
