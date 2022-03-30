import React from "react";
import type { CHIP_VARIANTS } from "@clutch-sh/core";
import { Chip, Grid, Link, Tooltip } from "@clutch-sh/core";

export interface ProjectInfoChip {
  text: string;
  title?: string;
  icon?: React.ReactElement;
  url?: string;
  variant?: typeof CHIP_VARIANTS[number];
}

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
