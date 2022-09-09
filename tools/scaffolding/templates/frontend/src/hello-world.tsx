import React from "react";

import type { WorkflowProps } from ".";

const HelloWorld: React.FC<WorkflowProps> = ({ heading }) => {
  return <h1>Hello World - {heading}</h1>;
};

export default HelloWorld;
