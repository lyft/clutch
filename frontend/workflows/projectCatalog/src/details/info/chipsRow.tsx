import React from "react";
import { Chip, Grid, Link, Tooltip } from "@clutch-sh/core";

import type { ProjectInfoChip } from "./types";

const ChipsRow = ({ chips = [] }: { chips: ProjectInfoChip[] }) => (
  <>
    {chips.map(({ variant = "neutral", text, icon, title, url }) => {
      const chipElem = <Chip variant={variant} label={text} size="small" icon={icon} />;
      return (
        <Grid item>
          <Tooltip title={title ?? text}>
            {url ? <Link href={url}>{chipElem}</Link> : chipElem}
          </Tooltip>
        </Grid>
      );
    })}
  </>
);

export default ChipsRow;
