import type { AxiosError } from "axios";

import type { ClientError, ClutchError, NetworkError } from "../errors";
import { client, errorInterceptor } from "../index";

describe("error interceptor", () => {
  const axiosError = {
    response: {
      data: {
        code: 5,
        message: "Could not find resource",
      },
    },
  } as AxiosError;

  describe("on axios error", () => {
    let err: Promise<ClientError | ClutchError | NetworkError>;
    beforeAll(() => {
      err = errorInterceptor({
        message: "Request timeout of 1ms reached",
        isAxiosError: true,
      } as AxiosError);
    });

    it("returns a ClientError", () => {
      return expect(err).rejects.toEqual({
        message: "Client Error: Request timeout of 1ms reached",
      });
    });
  });

  describe("on known error", () => {
    let err: Promise<ClientError | ClutchError | NetworkError>;
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
    let err: Promise<ClientError | ClutchError | NetworkError>;
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

    it("returns a NetworkError", () => {
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
