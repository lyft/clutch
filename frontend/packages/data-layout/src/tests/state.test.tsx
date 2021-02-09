import { ManagerAction, reducer } from "../state";

describe("Manager State", () => {
  describe("hydrate", () => {
    it("preserves existing data", () => {
      let state = {};
      state = reducer(state, {
        type: ManagerAction.SET,
        payload: { key: "layout1", value: { update: "value" } },
      });
      state = reducer(state, {
        type: ManagerAction.HYDRATE_END,
        payload: { key: "layout1", result: { hydrate: "value" } },
      });
      expect(state.layout1.data).toStrictEqual({ update: "value", hydrate: "value" });
    });
  });
});
