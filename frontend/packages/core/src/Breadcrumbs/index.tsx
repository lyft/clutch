import React from "react";
import { Breadcrumbs as MuiBreadcrumbs, Theme } from "@mui/material";
import { alpha } from "@mui/system";

import styled from "../styled";

import Breadcrumb from "./breadcrumb";
import type { BreadcrumbEntry } from "./types";

interface BreadcrumbsProps {
  entries: BreadcrumbEntry[];
}

const StyledBreadcrumbs = styled(MuiBreadcrumbs)(({ theme }: { theme: Theme }) => ({
  margin: theme.spacing(theme.clutch.spacing.sm, theme.clutch.spacing.none),
  "& .MuiBreadcrumbs-separator": {
    color: alpha(theme.colors.neutral[900], 0.6),
  },
}));

const Breadcrumbs = ({ entries }: BreadcrumbsProps) => (
  <StyledBreadcrumbs>
    {entries.map((entry: BreadcrumbEntry) => (
      <Breadcrumb key={entry.label} {...entry} />
    ))}
  </StyledBreadcrumbs>
);

export * from "./types";

export default Breadcrumbs;
