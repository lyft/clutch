import type { ProjectState } from "../types";
import { Group, isGlobalProjectState, isGroupState, isProjectState } from "../types";

describe("isGroupState", () => {
  it("returns false for undefined state", () => {
    const state = undefined;
    expect(isGroupState(state)).toBe(false);
  });

  it("matches group state types", () => {
    const state = { key: { checked: false } as ProjectState };
    expect(isGroupState(state)).toBe(true);
  });

  it("rejects other types", () => {
    const state = { key: {} as ProjectState };
    expect(isGroupState(state)).toBe(false);
  });
});

describe("isGlobalProjectState", () => {
  it("matches projects state types with only required fields", () => {
    const state = {
      [Group.PROJECTS]: { key: { checked: false } as ProjectState },
      [Group.UPSTREAM]: { key: { checked: false } as ProjectState },
      [Group.DOWNSTREAM]: { key: { checked: false } as ProjectState },
    };
    expect(isGlobalProjectState(state)).toBe(true);
  });

  it("matches projects state types with optional fields", () => {
    const state = {
      [Group.PROJECTS]: { key: { checked: false, custom: false } as ProjectState },
      [Group.UPSTREAM]: { key: { checked: false } as ProjectState },
      [Group.DOWNSTREAM]: { key: { checked: false } as ProjectState },
    };
    expect(isGlobalProjectState(state)).toBe(true);
  });

  it("rejects projects state types with incorrect types", () => {
    const state = {
      [Group.PROJECTS]: { key: { checked: "false" } },
      [Group.UPSTREAM]: { key: { checked: false } },
      [Group.DOWNSTREAM]: { key: { checked: false, custom: "false" } },
    };
    expect(isGlobalProjectState(state)).toBe(false);
  });

  it("rejects projects state types without required fields", () => {
    const state = {
      [Group.PROJECTS]: { key: { custom: true } as ProjectState },
      [Group.UPSTREAM]: { key: { custom: true } as ProjectState },
      [Group.DOWNSTREAM]: { key: { custom: false } as ProjectState },
    };
    expect(isGlobalProjectState(state)).toBe(false);
  });

  it("rejects other types", () => {
    const state = {
      [Group.PROJECTS]: { key: { checked: false } as ProjectState },
      [Group.UPSTREAM]: { key: {} as ProjectState },
    };
    expect(isGlobalProjectState(state)).toBe(false);
  });
});

describe("isProjectState", () => {
  it("matches project state types", () => {
    const state = { checked: false };
    expect(isProjectState(state)).toBe(true);
  });

  it("rejects other types", () => {
    const state = {};
    expect(isProjectState(state)).toBe(false);
  });
});
