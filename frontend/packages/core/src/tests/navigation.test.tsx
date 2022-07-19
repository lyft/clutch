import { convertSearchParam, useNavigate, useSearchParams } from "../navigation";

const mockNav = jest.fn();
const mockSearchParams = jest.fn();
const mockSetSearchParams = jest.fn();

jest.mock("react-router-dom", () => ({
  ...(jest.requireActual("react-router-dom") as any),
  useNavigate: () => mockNav,
  useSearchParams: () => [{ get: mockSearchParams }, mockSetSearchParams],
}));

describe("convertSearchParam", () => {
  it("handles empty objects", () => {
    const result = convertSearchParam(new URLSearchParams({}));
    expect(result).toBe("");
  });

  it("prepends ?", () => {
    const result = convertSearchParam(new URLSearchParams({ foo: "bar" }));
    expect(result).toBe("?foo=bar");
  });

  it("joins entries", () => {
    const result = convertSearchParam(new URLSearchParams({ foo: "bar", test: "42" }));
    expect(result).toBe("?foo=bar&test=42");
  });
});

describe("useNavigate", () => {
  const nav = useNavigate();

  beforeEach(() => {
    mockNav.mockReset();
  });

  describe("by default", () => {
    it("calls react-router navigate with partial path and no options", () => {
      nav("/foo");
      expect(mockNav).toHaveBeenCalledWith({ pathname: "/foo", search: "" }, {});
    });

    it("preserves only react router options", () => {
      nav("/foo", { utm: true, replace: true, state: { foo: "bar" } });
      expect(mockNav).toHaveBeenCalledWith(
        { pathname: "/foo", search: "" },
        { replace: true, state: { foo: "bar" } }
      );
    });
  });

  describe("with UTM search params", () => {
    beforeAll(() => {
      mockSearchParams.mockImplementation(
        k =>
          ({
            utm_source: "test",
            utm_medium: "run",
            randon: "value",
          }[k])
      );
    });

    afterAll(() => {
      mockSearchParams.mockImplementation();
    });

    it("preserves only UTM params from existing search", () => {
      nav("/foo");
      expect(mockNav).toHaveBeenCalledWith(
        { pathname: "/foo", search: "?utm_source=test&utm_medium=run" },
        {}
      );
    });

    it("merges UTM params with new search", () => {
      nav({
        pathname: "/foo",
        search: {
          new: "value",
        },
      });
      expect(mockNav).toHaveBeenCalledWith(
        { pathname: "/foo", search: "?new=value&utm_source=test&utm_medium=run" },
        {}
      );
    });

    it("can bypass UTM tracking", () => {
      nav("/foo", { utm: false });
      expect(mockNav).toHaveBeenCalledWith({ pathname: "/foo", search: "" }, {});
    });

    it("overwrites existing UTM params", () => {
      nav({
        pathname: "/foo",
        search: {
          utm_source: "newsource",
          utm_medium: "newmedium",
        },
      });
      expect(mockNav).toHaveBeenCalledWith(
        { pathname: "/foo", search: "?utm_source=newsource&utm_medium=newmedium" },
        {}
      );
    });
  });

  it("merges search param in path", () => {
    nav({
      pathname: "/foo?some=value",
      search: {
        search: "param",
      },
    });
    expect(mockNav).toHaveBeenCalledWith(
      { pathname: "/foo", search: "?search=param&some=value" },
      {}
    );
  });

  it("gives priority to pathname search param", () => {
    nav({
      pathname: "/foo?some=value",
      search: {
        some: "error",
      },
    });
    expect(mockNav).toHaveBeenCalledWith({ pathname: "/foo", search: "?some=value" }, {});
  });

  it("allows for custom partial path object", () => {
    nav({
      pathname: "/foo",
      search: { some: "param" },
    });
    expect(mockNav).toHaveBeenCalledWith(
      { hash: undefined, pathname: "/foo", search: "?some=param" },
      { replace: undefined, state: undefined }
    );
  });

  describe("origin", () => {
    const { location } = window;

    beforeAll(() => {
      delete window.location;
      // @ts-ignore
      window.location = { pathname: "/current/path", search: "?test=value" };
    });

    afterAll(() => {
      window.location = location;
    });

    it("can be persisted in state", () => {
      nav("/foo", { origin: true });
      expect(mockNav).toHaveBeenCalledWith(
        { pathname: "/foo", search: "" },
        { state: { origin: "/current/path?test=value" } }
      );
    });

    it("will not overwrite specified state origin", () => {
      nav("/foo", { origin: true, state: { origin: "/new/path" } });
      expect(mockNav).toHaveBeenCalledWith(
        { pathname: "/foo", search: "" },
        { state: { origin: "/new/path" } }
      );
    });
  });
});

describe("useSearchParams", () => {
  beforeAll(() => {
    mockSearchParams.mockImplementation(
      k =>
        ({
          foo: "bar",
        }[k])
    );
  });

  const [searchParams, setSearchParams] = useSearchParams();

  beforeEach(() => {
    mockSetSearchParams.mockReset();
  });

  it("passes through search params", () => {
    expect(searchParams.get("foo")).toBe("bar");
  });

  describe("setSearchParams", () => {
    it("calls react-router navigate with new search params and no options", () => {
      setSearchParams({ foo: "bar" });
      expect(mockSetSearchParams).toHaveBeenCalledWith({ foo: "bar" }, {});
    });

    it("preserves only react router options", () => {
      setSearchParams({ foo: "bar" }, { utm: true, replace: true, state: { foo: "bar" } });
      expect(mockSetSearchParams).toHaveBeenCalledWith(
        { foo: "bar" },
        { replace: true, state: { foo: "bar" } }
      );
    });
  });

  describe("with UTM search params", () => {
    beforeAll(() => {
      mockSearchParams.mockImplementation(
        k =>
          ({
            utm_source: "test",
            utm_medium: "run",
            randon: "value",
          }[k])
      );
    });

    afterAll(() => {
      mockSearchParams.mockImplementation();
    });

    it("preserves only UTM params from existing search", () => {
      setSearchParams({});
      expect(mockSetSearchParams).toHaveBeenCalledWith(
        { utm_source: "test", utm_medium: "run" },
        {}
      );
    });

    it("merges UTM params with new search", () => {
      setSearchParams({ new: "value" });
      expect(mockSetSearchParams).toHaveBeenCalledWith(
        { utm_source: "test", utm_medium: "run", new: "value" },
        {}
      );
    });

    it("can bypass UTM tracking", () => {
      setSearchParams({ foo: "bar" }, { utm: false });
      expect(mockSetSearchParams).toHaveBeenCalledWith({ foo: "bar" }, {});
    });

    it("overwrites existing UTM params", () => {
      setSearchParams({
        utm_source: "newsource",
        utm_medium: "newmedium",
      });
      expect(mockSetSearchParams).toHaveBeenCalledWith(
        {
          utm_source: "newsource",
          utm_medium: "newmedium",
        },
        {}
      );
    });
  });
});
