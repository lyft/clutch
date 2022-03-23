import React from "react";
import { Grid, styled, Typography } from "@clutch-sh/core";
import { Divider } from "@material-ui/core";

import { StyledLink } from "../card";

import type { AlertSummary } from "./types";

interface SummaryItemProps {
  count: number;
  title: string;
  color: string;
  url?: string;
}

const StyledDivider = styled(Divider)({
  color: "#A3A4B0",
  height: "24px",
  alignSelf: "center",
});

const SummaryItem = ({ count, title, color, url }: SummaryItemProps) => {
  const summ = (
    <Grid item>
      <Grid container item direction="column" alignItems="center">
        <Grid item>
          <Typography variant="subtitle2" color={color}>
            {count}
          </Typography>
        </Grid>
        <Grid item>
          <Typography variant="body4">{title}</Typography>
        </Grid>
      </Grid>
    </Grid>
  );

  return url ? <StyledLink href={url}>{summ}</StyledLink> : summ;
};

const SummaryRow = ({ open, triggered, acknowledged }: AlertSummary) => (
  <>
    {open && <SummaryItem count={open.count} title="Open" color="#CA4428" url={open.url} />}

    {open && triggered && <StyledDivider orientation="vertical" />}

    {triggered && (
      <SummaryItem count={triggered.count} title="Triggered" color="#3548D4" url={triggered.url} />
    )}

    {(open || triggered) && acknowledged && <StyledDivider orientation="vertical" />}

    {acknowledged && (
      <SummaryItem
        count={acknowledged.count}
        title="Acknowledged"
        color="#D97706"
        url={acknowledged.url}
      />
    )}
  </>
);

export default SummaryRow;
