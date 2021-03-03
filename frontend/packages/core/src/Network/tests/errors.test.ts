import type { AxiosError } from "axios";

import type { ClutchError, Help } from "../errors";
import clutchError from "../errors";
import * as grpcCodeToHttpStatus from "../grpc";

describe("clutch error", () => {
  const axiosError = {
    response: {
      data: {
        code: 5,
        message: "Could not find resource",
      },
    },
  } as AxiosError;

  describe("returns a basic ClutchError object", () => {
    let err: ClutchError;
    beforeAll(() => {
      err = clutchError(axiosError);
    });

    it("with a error code", () => {
      expect(err.code).toBe(5);
    });

    it("with a error messsage", () => {
      expect(err.message).toBe("Could not find resource");
    });

    it("with a status code", () => {
      expect(err.status.code).toBe(404);
    });

    it("with a status text", () => {
      expect(err.status.text).toBe("Not Found");
    });

    it("without details", () => {
      expect(err.details).toBeUndefined();
    });
  });

  describe("returns a detailed ClutchError object", () => {
    let err: ClutchError;
    beforeAll(() => {
      const complexAxiosError = { ...axiosError };
      complexAxiosError.response.data.details = [
        {
          "@type": "types.google.com/google.rpc.Help",
          links: [
            {
              description: "This is a link",
              url: "https://www.clutch.sh",
            },
          ],
        },
      ];
      err = clutchError(complexAxiosError);
    });

    it("with a list of details", () => {
      expect(err.details).toHaveLength(1);
    });

    it("with correct typing", () => {
      const helpDetails = err.details[0] as Help;
      expect(helpDetails.links).toHaveLength(1);
    });
  });

  it("converts code to an http status", () => {
    const grpcCodeToStatusMock = jest.spyOn(grpcCodeToHttpStatus, "default");
    clutchError(axiosError);
    expect(grpcCodeToStatusMock).toHaveBeenLastCalledWith(5);
  });
});
