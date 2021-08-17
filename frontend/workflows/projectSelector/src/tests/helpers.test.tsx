import {
  hydrateFromLocalState,
  isGroupState,
  isLocalState,
  isProjectState,
  LOCAL_STORAGE_STATE_KEY,
  writeToLocalState,
} from "../helpers";
import type { ProjectState, State } from "../types";
import { Group } from "../types";

describe("isGroupState", () => {
  it("returns false for undefined state", () => {
    const groupState = undefined;
    expect(isGroupState(groupState)).toBe(false);
  });

  it("matches group state types", () => {
    const groupState = { key: { checked: false } as ProjectState };
    expect(isGroupState(groupState)).toBe(true);
  });

  it("rejects other types", () => {
    const groupState = { key: {} as ProjectState };
    expect(isGroupState(groupState)).toBe(false);
  });
});

describe("isLocalState", () => {
  it("matches local state types", () => {
    const localState = {
      [Group.PROJECTS]: { key: { checked: false } as ProjectState },
      [Group.UPSTREAM]: { key: { checked: false } as ProjectState },
      [Group.DOWNSTREAM]: { key: { checked: false } as ProjectState },
    };
    expect(isLocalState(localState)).toBe(true);
  });

  it("rejects other types", () => {
    const localState = {
      [Group.PROJECTS]: { key: { checked: false } as ProjectState },
      [Group.UPSTREAM]: { key: {} as ProjectState },
    };
    expect(isLocalState(localState)).toBe(false);
  });
});

describe("isProjectState", () => {
  it("matches project state types", () => {
    const projectState = { checked: false };
    expect(isProjectState(projectState)).toBe(true);
  });

  it("rejects other types", () => {
    const projectState = {};
    expect(isProjectState(projectState)).toBe(false);
  });
});

describe("hydrateFromLocalState", () => {
  const localState = JSON.stringify({
    [Group.PROJECTS]: { a: { checked: false }, b: { checked: true } },
    [Group.UPSTREAM]: { b: { checked: true } },
    [Group.DOWNSTREAM]: { c: { checked: false } },
  });

  const state = {
    [Group.PROJECTS]: { a: { checked: true } },
    [Group.UPSTREAM]: {},
    [Group.DOWNSTREAM]: {},
    projectData: {},
    loading: false,
    error: undefined,
  } as State;

  beforeEach(() => {
    window.localStorage.clear();
  });

  it("returns existing state if local storage empty", () => {
    const finalState = hydrateFromLocalState(state);
    expect(finalState).toEqual(state);
  });

  it("removes local storage if invalid state", () => {
    window.localStorage.setItem(LOCAL_STORAGE_STATE_KEY, "{}");
    hydrateFromLocalState(state);
    expect(window.localStorage.getItem(LOCAL_STORAGE_STATE_KEY)).toBeNull();
  });

  it("returns existing state if invalid state", () => {
    window.localStorage.setItem(LOCAL_STORAGE_STATE_KEY, "{}");
    const finalState = hydrateFromLocalState(state);
    expect(finalState).toEqual(state);
  });

  it("returns existing state on any error", () => {
    window.localStorage.setItem(LOCAL_STORAGE_STATE_KEY, "foobar");
    const finalState = hydrateFromLocalState(state);
    expect(finalState).toEqual(state);
  });

  it("merges existing state with valid local state", () => {
    window.localStorage.setItem(LOCAL_STORAGE_STATE_KEY, localState);
    const finalState = hydrateFromLocalState(state);
    expect(finalState).toEqual({
      [Group.PROJECTS]: { a: { checked: false }, b: { checked: true } },
      [Group.UPSTREAM]: { b: { checked: true } },
      [Group.DOWNSTREAM]: { c: { checked: false } },
      projectData: {},
      loading: false,
      error: undefined,
    });
  });
});

describe("writeToLocalState", () => {
  const localState = {
    [Group.PROJECTS]: { a: { checked: false } },
    [Group.UPSTREAM]: { b: { checked: true } },
    [Group.DOWNSTREAM]: { c: { checked: false } },
  };
  const state = {
    ...localState,
    projectData: {},
    loading: false,
    error: undefined,
  } as State;

  beforeEach(() => {
    window.localStorage.clear();
  });

  it("writes stringified local state to local storage", () => {
    writeToLocalState(state);
    expect(window.localStorage.getItem(LOCAL_STORAGE_STATE_KEY)).toBe(JSON.stringify(localState));
  });
});
