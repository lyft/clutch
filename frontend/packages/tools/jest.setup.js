import Adapter from "@wojtekmaj/enzyme-adapter-react-17";
import Enzyme from "enzyme";

Enzyme.configure({ adapter: new Adapter() });

HTMLCanvasElement.prototype.getContext = () => {
  // return whatever getContext has to return
};
