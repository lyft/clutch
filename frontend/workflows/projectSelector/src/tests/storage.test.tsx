import {
  loadStoredState,
  LOCAL_STORAGE_STATE_KEY,
  storeState,
} from "../storage";
import type { State } from "../types";
import { Group } from "../types";

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
