import { grpcCodeToHttpCode, grpcCodeToText } from "../grpc";

describe("grpcCodeToHttpCode", () => {
  it("coverts known codes", () => {
    expect(grpcCodeToHttpCode(0)).toBe(200);
  });

  it("handles unknown codes", () => {
    expect(grpcCodeToHttpCode(-100)).toBe(500);
  });
});

describe("grpcCodeToText", () => {
  it("coverts known codes", () => {
    expect(grpcCodeToText(2)).toBe("Internal Server Error");
  });

  it("handles unknown codes", () => {
    expect(grpcCodeToText(-100)).toBe("Unknown Code");
  });
});
