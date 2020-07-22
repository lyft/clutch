import React from "react";
import { ExpansionPanel, TreeTable } from "@clutch-sh/core";
import _ from "lodash";

const Stats = ({ stats }) => {
  const structuredStats = {};
  stats.stats.forEach(stat => {
    if (stat.value > 0) {
      _.setWith(structuredStats, stat.key, stat.value, Object);
    }
  });

  const status = `(${stats.stats.length} total)`;
  return (
    <ExpansionPanel heading="Stats" summary={status}>
      <TreeTable data={structuredStats} />
    </ExpansionPanel>
  );
};

export default Stats;
