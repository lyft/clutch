import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";

interface ContextProps {
  projectInfo: IClutch.core.project.v1.IProject | null;
  projectId: string | null;
}

const ProjectDetailsContext = React.createContext<ContextProps | undefined>(undefined);

const useProjectDetailsContext = () => {
  return React.useContext<ContextProps | undefined>(ProjectDetailsContext);
};

export { ProjectDetailsContext, useProjectDetailsContext };
