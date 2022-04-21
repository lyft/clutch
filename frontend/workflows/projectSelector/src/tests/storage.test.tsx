import { COMPONENT_NAME, getLocalState, loadStoredState, STORAGE_STATE_KEY } from "../storage";
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
    projectErrors: undefined,
    loading: false,
    error: undefined,
  } as State;

  const retrieveLocalData = (key: string) => {
    const localData = window.localStorage.getItem(key);

    if (localData) {
      try {
        return JSON.parse(localData);
      } catch (_) {
        return localData;
      }
    }
  };

  const removeLocalData = (key: string) => window.localStorage.removeItem(key);

  const retrieveDataMock = (name: string, key: string, defaultData: State) =>
    retrieveLocalData(key) ?? defaultData;

  const removeDataMock = (name: string, key: string, local: boolean) => removeLocalData(key);

  beforeEach(() => {
    window.localStorage.clear();
  });

  it("returns existing state if local storage empty", () => {
    const finalState = loadStoredState(state, retrieveDataMock, removeDataMock);
    expect(finalState).toEqual(state);
  });

  it("removes local storage if invalid state", () => {
    window.localStorage.setItem(STORAGE_STATE_KEY, "{}");
    loadStoredState(state, retrieveDataMock, removeDataMock);
    expect(window.localStorage.getItem(STORAGE_STATE_KEY)).toBeNull();
  });

  it("returns existing state if invalid state", () => {
    window.localStorage.setItem(STORAGE_STATE_KEY, "{}");
    const finalState = loadStoredState(state, retrieveDataMock, removeDataMock);
    expect(finalState).toEqual(state);
  });

  it("returns existing state on any error", () => {
    window.localStorage.setItem(STORAGE_STATE_KEY, "foobar");
    const finalState = loadStoredState(state, retrieveDataMock, removeDataMock);
    expect(finalState).toEqual(state);
  });

  it("merges existing state with valid local state", () => {
    window.localStorage.setItem(STORAGE_STATE_KEY, storedState);
    const finalState = loadStoredState(state, retrieveDataMock, removeDataMock);
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
  const globalState = {
    [Group.PROJECTS]: { a: { checked: false } },
    [Group.UPSTREAM]: { b: { checked: true } },
    [Group.DOWNSTREAM]: { c: { checked: false } },
  };
  const state = {
    ...globalState,
    projectData: {},
    projectErrors: undefined,
    loading: false,
    error: undefined,
  } as State;

  const storeLocalData = (key: string, data: unknown) => {
    try {
      window.localStorage.setItem(key, JSON.stringify(data));
    } catch (e) {
      // eslint-disable-next-line no-console
      console.error("Error saving to local storage", e);
    }
  };

  const storeDataMock = (name: string, key: string, data: unknown, local: boolean) =>
    storeLocalData(key, data);

  beforeEach(() => {
    window.localStorage.clear();
  });

  it("writes stringified projects state to local storage", () => {
    storeDataMock(COMPONENT_NAME, STORAGE_STATE_KEY, getLocalState(state), true);
    expect(window.localStorage.getItem(STORAGE_STATE_KEY)).toBe(JSON.stringify(globalState));
  });
});
