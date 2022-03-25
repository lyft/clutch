import React from "react";

import type { ProjectInfo } from "./info/types";

interface ContextProps {
  projectInfo: ProjectInfo | null;
}

const ProjectDetailsContext = React.createContext<ContextProps | undefined>(undefined);

const useProjectDetailsContext = () => {
  return React.useContext<ContextProps | undefined>(ProjectDetailsContext);
};

export { ProjectDetailsContext, useProjectDetailsContext };
