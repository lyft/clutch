import React from "react";
import type { CHIP_VARIANTS } from "@clutch-sh/core";
import { Chip, Grid, Link, Tooltip } from "@clutch-sh/core";
import ChevronRightIcon from "@mui/icons-material/ChevronRight";
import type { LinkProps } from "@mui/material";

export interface ProjectInfoChip {
  text: string;
  title?: string;
  icon?: React.ReactElement;
  url?: string;
  urlTarget?: LinkProps["target"];
  variant?: typeof CHIP_VARIANTS[number];
}

const ChipsRow = ({ chips = [] }: { chips: ProjectInfoChip[] }) => (
  <>
    {chips.map(({ variant = "neutral", text, icon, title, url, urlTarget }) => {
      const chipText = (
        <Grid container direction="row" wrap="nowrap">
          {text}
          {url && <ChevronRightIcon fontSize="small" style={{ marginRight: "-12px" }} />}
        </Grid>
      );
      const chipElem = <Chip variant={variant} label={chipText} size="small" icon={icon} />;

      if (!urlTarget) {
        const externalRoute = url && url.startsWith("http");

        // eslint-disable-next-line no-param-reassign
        urlTarget = externalRoute ? "_blank" : "_self";
      }

      return (
        <Tooltip title={title ?? text} key={`chip-${title}`}>
          <Grid item>
            {url ? (
              <Link href={url} {...(urlTarget ? { target: urlTarget } : {})}>
                {chipElem}
              </Link>
            ) : (
              chipElem
            )}
          </Grid>
        </Tooltip>
      );
    })}
  </>
);

export default ChipsRow;
