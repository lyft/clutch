import React from "react";
import { ExpansionPanel, TreeTable } from "@clutch-sh/core";
import _ from "lodash";

const Runtime = ({ runtime }) => {
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
