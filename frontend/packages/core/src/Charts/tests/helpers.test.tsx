import {
  calculateDomainEdges,
  calculateTicks,
  dateTimeFormatter,
  getLeftSideAndIntervalForTicks,
  getMinAndMaxOfRangeUsingKey,
  isoTimeFormatter,
} from "../helpers";

describe("helper functions for charts", () => {
  describe("isoTimeFormatter", () => {
    test("a timestamp", () => {
      expect(isoTimeFormatter(60)).toBe("1970-01-01T00:00:00.060Z");
    });
  });

  describe("dateTimeFormatter", () => {
    test("a timestamp", () => {
      expect(dateTimeFormatter(60)).toBe("Wed Dec 31 1969");
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
  describe("getMinAndMaxOfRangeUsingKey", () => {
    test("the key exists in the data", () => {
      expect(getMinAndMaxOfRangeUsingKey(testDataMeowWoof, "meow")).toStrictEqual({
        min: 1,
        max: 100,
      });
    });

    test("the key does not exist in the data", () => {
      expect(getMinAndMaxOfRangeUsingKey(testDataMeowWoof, "foo")).toStrictEqual({
        min: Infinity,
        max: -Infinity,
      });
    });

    test("data is null", () => {
      expect(getMinAndMaxOfRangeUsingKey(null, "meow")).toStrictEqual({ min: null, max: null });
    });

    test("data is undefined", () => {
      expect(getMinAndMaxOfRangeUsingKey(undefined, "meow")).toStrictEqual({
        min: null,
        max: null,
      });
    });

    test("data is empty", () => {
      expect(getMinAndMaxOfRangeUsingKey([], "meow")).toStrictEqual({
        min: Infinity,
        max: -Infinity,
      });
    });

    test("min and max are the same", () => {
      expect(getMinAndMaxOfRangeUsingKey(testDataMeowWoof, "woof")).toStrictEqual({
        min: 2,
        max: 2,
      });
    });
  });

  describe("calculateDomainEdges", () => {
    test("edgeRatio is negative", () => {
      expect(calculateDomainEdges(testDataMeowWoof, "meow", -1)).toStrictEqual([1, 100]);
    });

    test("edgeRation is 0", () => {
      expect(calculateDomainEdges(testDataMeowWoof, "meow", 0)).toStrictEqual([1, 100]);
    });

    test("happy path", () => {
      expect(calculateDomainEdges(testDataMeowWoof, "meow", 0.1)).toStrictEqual([-8.9, 109.9]);
    });

    test("dataKey does not exist in data", () => {
      expect(calculateDomainEdges(testDataMeowWoof, "foo", 0.1)).toStrictEqual([
        Infinity,
        -Infinity,
      ]);
    });

    test("data is null", () => {
      expect(calculateDomainEdges(null, "meow", 0.1)).toStrictEqual([null, null]);
    });

    test("data is undefined", () => {
      expect(calculateDomainEdges(undefined, "meow", 0.1)).toStrictEqual([null, null]);
    });

    test("data is empty", () => {
      expect(calculateDomainEdges([], "meow", 0.1)).toStrictEqual([Infinity, -Infinity]);
    });

    test("max and min are equal in data", () => {
      expect(calculateDomainEdges(testDataMeowWoof, "woof", 0.1)).toStrictEqual([1.8, 2.2]);
    });
  });

  describe("getLeftSideAndIntervalForTicks", () => {
    test("min is greater than max", () => {
      expect(getLeftSideAndIntervalForTicks(1, 0)).toStrictEqual({
        leftSide: null,
        interval: null,
      });
    });

    // Note that 15000 refers to milliseconds (15 seconds), the smallest interval possible
    test("max is equal to min", () => {
      expect(getLeftSideAndIntervalForTicks(1, 1)).toStrictEqual({ leftSide: 0, interval: 15000 });
    });

    test("happy path", () => {
      expect(getLeftSideAndIntervalForTicks(1, 2)).toStrictEqual({ leftSide: 0, interval: 15000 });
    });

    test("min is negative", () => {
      expect(getLeftSideAndIntervalForTicks(-1, 2)).toStrictEqual({ leftSide: 0, interval: 15000 });
    });

    test("min is 0", () => {
      expect(getLeftSideAndIntervalForTicks(0, 2)).toStrictEqual({ leftSide: 0, interval: 15000 });
    });

    test("happy path 2", () => {
      expect(getLeftSideAndIntervalForTicks(1, 3)).toStrictEqual({ leftSide: 0, interval: 15000 });
    });
  });

  describe("calculateTicks", () => {
    test("data is null", () => {
      expect(calculateTicks(null, "meow")).toStrictEqual([]);
    });

    test("data is undefined", () => {
      expect(calculateTicks(undefined, "meow")).toStrictEqual([]);
    });

    test("data is empty", () => {
      expect(calculateTicks([], "meow")).toStrictEqual([]);
    });

    test("happy path", () => {
      expect(calculateTicks(testDataMeowWoof, "meow")).toStrictEqual([0]);
    });
  });
});
