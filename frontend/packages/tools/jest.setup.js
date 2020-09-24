import Enzyme from "enzyme";
import Adapter from "enzyme-adapter-react-16";

Enzyme.configure({ adapter: new Adapter() });

HTMLCanvasElement.prototype.getContext = () => { 
  // return whatever getContext has to return
};