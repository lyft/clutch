import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { TreeTable } from "@clutch-sh/core";
import _ from "lodash";

interface RuntimeProps {
  runtime: IClutch.envoytriage.v1.IRuntime;
}

const Runtime: React.FC<RuntimeProps> = ({ runtime }) => {
  const structuredEntries = {};
  runtime.entries.forEach(entry => {
    _.setWith(structuredEntries, entry.key, entry.value, Object);
  });

  return <TreeTable data={structuredEntries} />;
};

export default Runtime;
