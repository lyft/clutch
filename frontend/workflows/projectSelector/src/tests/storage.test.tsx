import {
  isGroupState,
  isProjectsState,
  isProjectState,
  loadStoredState,
  LOCAL_STORAGE_STATE_KEY,
  storeState,
} from "../storage";
import type { ProjectState, State } from "../types";
import { Group } from "../types";

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

describe("isProjectsState", () => {
  it("matches projects state types", () => {
    const state = {
      [Group.PROJECTS]: { key: { checked: false } as ProjectState },
      [Group.UPSTREAM]: { key: { checked: false } as ProjectState },
      [Group.DOWNSTREAM]: { key: { checked: false } as ProjectState },
    };
    expect(isProjectsState(state)).toBe(true);
  });

  it("rejects other types", () => {
    const state = {
      [Group.PROJECTS]: { key: { checked: false } as ProjectState },
      [Group.UPSTREAM]: { key: {} as ProjectState },
    };
    expect(isProjectsState(state)).toBe(false);
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

describe("loadStoredState", () => {
  const storedState = JSON.stringify({
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
    const finalState = loadStoredState(state);
    expect(finalState).toEqual(state);
  });

  it("removes local storage if invalid state", () => {
    window.localStorage.setItem(LOCAL_STORAGE_STATE_KEY, "{}");
    loadStoredState(state);
    expect(window.localStorage.getItem(LOCAL_STORAGE_STATE_KEY)).toBeNull();
  });

  it("returns existing state if invalid state", () => {
    window.localStorage.setItem(LOCAL_STORAGE_STATE_KEY, "{}");
    const finalState = loadStoredState(state);
    expect(finalState).toEqual(state);
  });

  it("returns existing state on any error", () => {
    window.localStorage.setItem(LOCAL_STORAGE_STATE_KEY, "foobar");
    const finalState = loadStoredState(state);
    expect(finalState).toEqual(state);
  });

  it("merges existing state with valid local state", () => {
    window.localStorage.setItem(LOCAL_STORAGE_STATE_KEY, storedState);
    const finalState = loadStoredState(state);
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

describe("storeState", () => {
  const projectsState = {
    [Group.PROJECTS]: { a: { checked: false } },
    [Group.UPSTREAM]: { b: { checked: true } },
    [Group.DOWNSTREAM]: { c: { checked: false } },
  };
  const state = {
    ...projectsState,
    projectData: {},
    loading: false,
    error: undefined,
  } as State;

  beforeEach(() => {
    window.localStorage.clear();
  });

  it("writes stringified projects state to local storage", () => {
    storeState(state);
    expect(window.localStorage.getItem(LOCAL_STORAGE_STATE_KEY)).toBe(
      JSON.stringify(projectsState)
    );
  });
});
