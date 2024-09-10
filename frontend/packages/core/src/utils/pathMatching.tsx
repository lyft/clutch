import { matchPath } from "react-router";

const findPathMatchList = (locationPathname: string, pathstoMatch: string[]) => {
  let pathFound = false;

  pathstoMatch?.forEach((path: string) => {
    const match = matchPath({ path }, locationPathname);

    if (match) {
      pathFound = true;
    }
  });

  return pathFound;
};

export default findPathMatchList;
