import { matchPath } from "react-router";

const checkPathMatchList = (locationPathname: string, pathstoMatch: string[]) => {
  pathstoMatch.forEach(path => {
    const match = matchPath({ path }, locationPathname);

    console.log("match: ", match);
  });
};

export default checkPathMatchList;
