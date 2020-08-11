import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { ExpansionPanel, TreeTable } from "@clutch-sh/core";
import _ from "lodash";

interface RuntimeProps {
  runtime: IClutch.envoytriage.v1.IRuntime;
}

const Runtime: React.FC<RuntimeProps> = ({ runtime }) => {
  const structuredEntries = {};
  let status = "";
  runtime.entries.forEach(entry => {
    _.setWith(structuredEntries, entry.key, entry.value, Object);
  });

  status = `(${runtime.entries.length} total)`;
  return (
    <ExpansionPanel heading="Runtime" summary={status}>
      <TreeTable data={structuredEntries} />
    </ExpansionPanel>
  );
};

export default Runtime;
