import type { AxiosError, AxiosResponse } from "axios";

import type { ClutchError } from "../errors";
import { client, errorInterceptor, successInterceptor } from "../index";

describe("success interceptor", () => {
  const { location } = window;

  beforeAll(() => {
    delete window.location;
  });

  afterAll(() => {
    window.location = location;
  });

  describe("on auth url response data", () => {
    let response: () => AxiosResponse;
    beforeEach(() => {
      response = () =>
        successInterceptor({
          data: {
            authUrl: "https://clutch.sh/auth",
          },
        } as AxiosResponse);
    });

    it("redirects to provided url", () => {
      expect(() => response()).toThrow();
      expect(window.location).toBe("https://clutch.sh/auth");
    });

    it("throws a ClutchError", () => {
      expect(() => response()).toThrow({
        message: "Authentication Expired",
        status: {
          code: 401,
          text: "Authentication Expired",
        },
      } as ClutchError);
    });
  });
});

describe("error interceptor", () => {
  describe("on axios error", () => {
    let err: Promise<ClutchError>;
    beforeAll(() => {
      err = errorInterceptor({
        message: "Request timeout of 1ms reached",
        isAxiosError: true,
      } as AxiosError);
    });

    it("returns a ClutchError", () => {
      return expect(err).rejects.toEqual({
        status: {
          code: 500,
          text: "Client Error",
        },
        message: "Request timeout of 1ms reached",
      });
    });
  });

  describe("on known error", () => {
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
    let err: Promise<ClutchError>;
    beforeAll(() => {
      err = errorInterceptor(axiosError);
    });

    it("returns a ClutchError", () => {
      return expect(err).rejects.toEqual({
        code: 5,
        message: "Could not find resource",
        status: {
          code: 404,
          text: "Not Found",
        },
      });
    });
  });

  describe("on unknown error", () => {
    let err: Promise<ClutchError>;
    beforeAll(() => {
      err = errorInterceptor({
        isAxiosError: false,
        message: "Unauthorized to perform action",
        response: {
          data: {},
          status: 401,
          statusText: "Unauthenticated",
          headers: {},
          config: {},
        },
        config: {},
        name: "foobar",
        toJSON: () => {
          return {};
        },
      });
    });

    it("returns a ClutchError", () => {
      return expect(err).rejects.toEqual({
        data: {},
        message: "Unauthorized to perform action",
        status: {
          code: 401,
          text: "Unauthenticated",
        },
      });
    });
  });
});

describe("axios client", () => {
  it("treats status codes >= 400 as error", () => {
    expect(client.defaults.validateStatus(400)).toBe(false);
  });

  it("treats status codes >= 500 as error", () => {
    expect(client.defaults.validateStatus(500)).toBe(false);
  });

  it("treats status codes < 400 as success", () => {
    expect(client.defaults.validateStatus(399)).toBe(true);
  });
});
