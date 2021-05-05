import Enzyme from "enzyme";
import Adapter from "@wojtekmaj/enzyme-adapter-react-17";

Enzyme.configure({ adapter: new Adapter() });

HTMLCanvasElement.prototype.getContext = () => { 
  // return whatever getContext has to return
};