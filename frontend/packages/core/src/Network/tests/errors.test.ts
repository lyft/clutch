import type { AxiosError } from "axios";

import type { ClutchError, ClutchErrorDetails, Help } from "../errors";
import { grpcResponseToError, isClutchErrorDetails, isHelpDetails } from "../errors";

describe("clutch error", () => {
  const axiosError = {
    response: {
      status: 404,
      statusText: "Not Found",
      data: {
        code: 5,
        message: "Could not find resource",
      },
    },
  } as AxiosError;

  describe("returns a basic ClutchError object", () => {
    let err: ClutchError;
    beforeAll(() => {
      err = grpcResponseToError(axiosError);
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
          "@type": "types.googleapis.com/google.rpc.Help",
          links: [
            {
              description: "This is a link",
              url: "https://www.clutch.sh",
            },
          ],
        },
      ];
      err = grpcResponseToError(complexAxiosError);
    });

    it("with a list of details", () => {
      expect(err.details).toHaveLength(1);
    });

    it("with correct typing", () => {
      const helpDetails = err.details[0] as Help;
      expect(helpDetails.links).toHaveLength(1);
    });
  });
});

describe("isHelpDetails", () => {
  it("returns true for help details", () => {
    const details = {
      _type: "types.googleapis.com/google.rpc.Help",
      links: [
        {
          description: "Please file a ticket here for more help.",
          url: "https://www.example.com",
        },
      ],
    } as Help;

    expect(isHelpDetails(details)).toBe(true);
  });

  it("returns false for non-help details", () => {
    const details = {
      type: "unknownType",
      something: [
        {
          key: "value",
        },
      ],
    };
    expect(isHelpDetails(details)).toBe(false);
  });
});

describe("isClutchErrorDetails", () => {
  it("returns true for Clutch specific error details", () => {
    const details = {
      _type: "type.googleapis.com/clutch.api.v1.ErrorDetails",
      wrapped: [
        {
          code: 2,
          message: "core-staging-0: yikes",
        },
        {
          code: 16,
          message: "core-staging-1: nono",
        },
      ],
    } as ClutchErrorDetails;

    expect(isClutchErrorDetails(details)).toBe(true);
  });

  it("returns false for non-Clutch specific error details", () => {
    const details = {
      type: "unknownType",
      something: [
        {
          key: "value",
        },
      ],
    };
    expect(isClutchErrorDetails(details)).toBe(false);
  });
});
