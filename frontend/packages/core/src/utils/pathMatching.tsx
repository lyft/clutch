import { matchPath } from "react-router";

const findPathMatchList = (locationPathname: string, pathsToMatch: string[]) => {
  const pathFound = pathsToMatch?.find((path: string) => matchPath({ path }, locationPathname));

  return pathFound;
};

export default findPathMatchList;
