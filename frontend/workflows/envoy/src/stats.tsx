import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { ExpansionPanel, TreeTable } from "@clutch-sh/core";
import _ from "lodash";

interface StatsProps {
  stats: IClutch.envoytriage.v1.IStats;
}

const Stats: React.FC<StatsProps> = ({ stats }) => {
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
