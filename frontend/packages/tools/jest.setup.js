import Adapter from "@wojtekmaj/enzyme-adapter-react-17";
import Enzyme from "enzyme";

Enzyme.configure({ adapter: new Adapter() });

HTMLCanvasElement.prototype.getContext = () => {
  // return whatever getContext has to return
};

const localStorageMock = (() => {
  let store = {};

  return {
    getItem: key => {
      return store[key] || null;
    },
    setItem: (key, value) => {
      store[key] = value.toString();
    },
    removeItem: key => {
      delete store[key];
    },
    clear: () => {
      store = {};
    },
  };
})();

Object.defineProperty(window, "localStorage", {
  value: localStorageMock,
  writable: true,
});
