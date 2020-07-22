import { grpcCodeToHttpStatus, parseErrorMessage } from "../network";

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

describe("parseErrorMessage", () => {
  describe("with error missing newlines", () => {
    let errorDetails;
    const errorMsg = "This is an error";
    beforeAll(() => {
      errorDetails = parseErrorMessage(errorMsg);
    });

    it("has populated summary", () => {
      expect(errorDetails.summary).toEqual(errorMsg);
    });

    it("has empty details", () => {
      expect(errorDetails.details).toEqual("");
    });
  });

  describe("with error containing newlines", () => {
    let errorDetails;
    const summary = "This is an error";
    const details = "With additional content";
    const errorMsg = `${summary}\n${details}`;
    beforeAll(() => {
      errorDetails = parseErrorMessage(errorMsg);
    });
    it("has populated summary", () => {
      expect(errorDetails.summary).toEqual(summary);
    });

    it("has populated details", () => {
      expect(errorDetails.details).toEqual(details);
    });
  });
});
