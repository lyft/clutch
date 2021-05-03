import { ManagerAction, reducer } from "../state";

describe("Manager State", () => {
  describe("hydrate end", () => {
    describe("preserves existing data", () => {
      it("of lists", () => {
        let state = {};
        state = reducer(state, {
          type: ManagerAction.SET,
          payload: { key: "layout", value: [{ update: "value" }] },
        });
        state = reducer(state, {
          type: ManagerAction.HYDRATE_END,
          payload: { key: "layout", result: [{ hydrate: "value" }] },
        });
        expect(state.layout.data).toEqual([{ update: "value" }, { hydrate: "value" }]);
      });

      it("of objects", () => {
        let state = {};
        state = reducer(state, {
          type: ManagerAction.SET,
          payload: { key: "layout", value: { update: "value" } },
        });
        state = reducer(state, {
          type: ManagerAction.HYDRATE_END,
          payload: { key: "layout", result: { hydrate: "value" } },
        });
        expect(state.layout.data).toStrictEqual({ update: "value", hydrate: "value" });
      });
    });

    it("gives priority to new data", () => {
      let state = {};
      state = reducer(state, {
        type: ManagerAction.SET,
        payload: { key: "layout", value: { update: "initialValue" } },
      });
      state = reducer(state, {
        type: ManagerAction.HYDRATE_END,
        payload: { key: "layout", result: { update: "endingValue" } },
      });
      expect(state.layout.data).toStrictEqual({ update: "endingValue" });
    });

    describe("on type mismatch", () => {
      describe("overwrites existing data", () => {
        it("of objects", () => {
          let state = {};
          state = reducer(state, {
            type: ManagerAction.SET,
            payload: { key: "layout", value: { update: "value" } },
          });
          state = reducer(state, {
            type: ManagerAction.HYDRATE_END,
            payload: { key: "layout", result: [{ hydrate: "value" }] },
          });
          expect(state.layout.data).toStrictEqual([{ hydrate: "value" }]);
        });

        it("of lists", () => {
          let state = {};
          state = reducer(state, {
            type: ManagerAction.SET,
            payload: { key: "layout", value: [{ update: "value" }] },
          });
          state = reducer(state, {
            type: ManagerAction.HYDRATE_END,
            payload: { key: "layout", result: { hydrate: "value" } },
          });
          expect(state.layout.data).toStrictEqual({ hydrate: "value" });
        });
      });
    });
  });
});
