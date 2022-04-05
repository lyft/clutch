import React from "react";
import type { CHIP_VARIANTS } from "@clutch-sh/core";
import { Chip, Grid, Link, Tooltip } from "@clutch-sh/core";
import ChevronRightIcon from "@material-ui/icons/ChevronRight";

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
      const chipText = (
        <Grid container direction="row" wrap="nowrap">
          {text}
          {url && <ChevronRightIcon fontSize="small" style={{ marginRight: "-12px" }} />}
        </Grid>
      );
      const chipElem = <Chip variant={variant} label={chipText} size="small" icon={icon} />;
      return (
        <Tooltip title={title ?? text}>
          <Grid item>{url ? <Link href={url}>{chipElem}</Link> : chipElem}</Grid>
        </Tooltip>
      );
    })}
  </>
);

export default ChipsRow;
