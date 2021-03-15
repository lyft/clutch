import { grpcCodeToHttpCode, grpcCodeToText } from "../grpc";

describe("grpcCodeToHttpCode", () => {
  it("coverts known codes", () => {
    expect(grpcCodeToHttpCode(0)).toEqual(200);
  });

  it("handles unknown codes", () => {
    expect(grpcCodeToHttpCode(-100)).toEqual(500);
  });
});

describe("grpcCodeToText", () => {
  it("coverts known codes", () => {
    expect(grpcCodeToText(2)).toEqual("Internal Server Error");
  });

  it("handles unknown codes", () => {
    expect(grpcCodeToText(-100)).toEqual("Unknown Code");
  });
});
