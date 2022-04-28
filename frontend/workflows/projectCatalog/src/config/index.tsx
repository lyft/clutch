import React from "react";
import { useParams } from "react-router-dom";
import type { ProjectDetailsConfigWorkflowProps } from "..";

const Config: React.FC<ProjectDetailsConfigWorkflowProps> = ({ children, defaultRoute }) => {
  const { projectId, configType = defaultRoute } = useParams();

  React.useEffect(() => {
    console.log("I have chitlins");
  }, [children]);

  return (
    <>
      <div>Hello World!</div>
    </>
  );
};

export default Config;
