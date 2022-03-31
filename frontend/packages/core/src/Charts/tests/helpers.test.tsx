import {
  calculateDomainEdges,
  calculateTicks,
  getLeftSideAndIntervalForTicks,
  getMinAndMaxOfRangeUsingKey,
  isoTimeFormatter,
} from "../helpers";

describe("helper functions for charts", () => {
  describe("isoTimeFormatter()", () => {
    test("it returns a timestamp", () => {
      expect(isoTimeFormatter(60)).toBe("1970-01-01T00:00:00.060Z");
    });
  });

  const testDataMeowWoof = [
    {
      meow: "1",
      woof: "2",
    },
    {
      meow: "100",
      woof: "2",
    },
    {
      meow: "50",
      woof: "2",
    },
  ];
  describe("getMinAndMaxOfRangeUsingKey()", () => {
    test("it returns a min and max with proper input data", () => {
      expect(getMinAndMaxOfRangeUsingKey(testDataMeowWoof, "meow")).toStrictEqual({
        min: 1,
        max: 100,
      });
    });

    test("it returns non finite numbers with improper input data", () => {
      expect(getMinAndMaxOfRangeUsingKey(testDataMeowWoof, "foo")).toStrictEqual({
        min: Infinity,
        max: -Infinity,
      });
    });

    test("it returns null vals with null input data", () => {
      expect(getMinAndMaxOfRangeUsingKey(null, "meow")).toStrictEqual({ min: null, max: null });
    });

    test("it returns null vals with undefined input data", () => {
      expect(getMinAndMaxOfRangeUsingKey(undefined, "meow")).toStrictEqual({
        min: null,
        max: null,
      });
    });

    test("it returns non finite numbers with empty input data", () => {
      expect(getMinAndMaxOfRangeUsingKey([], "meow")).toStrictEqual({
        min: Infinity,
        max: -Infinity,
      });
    });

    test("it returns the min and max as the same when the input data has that", () => {
      expect(getMinAndMaxOfRangeUsingKey(testDataMeowWoof, "woof")).toStrictEqual({
        min: 2,
        max: 2,
      });
    });
  });

  describe("calculateDomainEdges()", () => {
    test("it works when edge ratio is negative", () => {
      expect(calculateDomainEdges(testDataMeowWoof, "meow", -1)).toStrictEqual([1, 100]);
    });

    test("it works when edge ratio is 0", () => {
      expect(calculateDomainEdges(testDataMeowWoof, "meow", 0)).toStrictEqual([1, 100]);
    });

    test("it works when edge ratio is .1 (happy path)", () => {
      expect(calculateDomainEdges(testDataMeowWoof, "meow", 0.1)).toStrictEqual([-8.9, 109.9]);
    });

    test("it returns non finite vals when the key does not exist in the input data", () => {
      expect(calculateDomainEdges(testDataMeowWoof, "foo", 0.1)).toStrictEqual([
        Infinity,
        -Infinity,
      ]);
    });

    test("it returns null vals when input is null", () => {
      expect(calculateDomainEdges(null, "meow", 0.1)).toStrictEqual([null, null]);
    });

    test("it returns null vals when input is undefined", () => {
      expect(calculateDomainEdges(undefined, "meow", 0.1)).toStrictEqual([null, null]);
    });

    test("it returns non finite vals when input is empty", () => {
      expect(calculateDomainEdges([], "meow", 0.1)).toStrictEqual([Infinity, -Infinity]);
    });

    test("it works (happy path 2)", () => {
      expect(calculateDomainEdges(testDataMeowWoof, "woof", 0.1)).toStrictEqual([1.8, 2.2]);
    });
  });

  describe("getLeftSideAndIntervalForTicks()", () => {
    test("it returns null when min is greater than max", () => {
      expect(getLeftSideAndIntervalForTicks(1, 0)).toStrictEqual({
        leftSide: null,
        interval: null,
      });
    });

    // Note that 15000 refers to milliseconds (15 seconds), the smallest interval possible
    test("it works fine when max and min are equal", () => {
      expect(getLeftSideAndIntervalForTicks(1, 1)).toStrictEqual({ leftSide: 0, interval: 15000 });
    });

    test("the happy path works fine", () => {
      expect(getLeftSideAndIntervalForTicks(1, 2)).toStrictEqual({ leftSide: 0, interval: 15000 });
    });

    test("it works when min is negative", () => {
      expect(getLeftSideAndIntervalForTicks(-1, 2)).toStrictEqual({ leftSide: 0, interval: 15000 });
    });

    test("it works when min is 0", () => {
      expect(getLeftSideAndIntervalForTicks(0, 2)).toStrictEqual({ leftSide: 0, interval: 15000 });
    });

    test("a separate exercise of the happy path", () => {
      expect(getLeftSideAndIntervalForTicks(1, 3)).toStrictEqual({ leftSide: 0, interval: 15000 });
    });
  });

  describe("calculateTicks()", () => {
    test("it returns empty when input is null", () => {
      expect(calculateTicks(null, "meow")).toStrictEqual([]);
    });

    test("it returns empty when input is undefined", () => {
      expect(calculateTicks(undefined, "meow")).toStrictEqual([]);
    });

    test("it returns empty when input is empty", () => {
      expect(calculateTicks([], "meow")).toStrictEqual([]);
    });

    test("it returns a single [0] when input is normal", () => {
      expect(calculateTicks(testDataMeowWoof, "meow")).toStrictEqual([0]);
    });
  });
});
