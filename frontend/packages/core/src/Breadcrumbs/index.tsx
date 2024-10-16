import React from "react";
import { Breadcrumbs as MuiBreadcrumbs, Theme } from "@mui/material";
import { alpha } from "@mui/system";

import styled from "../styled";

import Breadcrumb from "./breadcrumb";

export interface BreadcrumbEntry {
  label: string;
  url?: string;
}

export interface BreadcrumbsProps {
  entries: BreadcrumbEntry[];
}

const StyledBreadcrumbs = styled(MuiBreadcrumbs)(({ theme }: { theme: Theme }) => ({
  margin: "8px 0px",
  "& .MuiBreadcrumbs-separator": {
    color: alpha(theme.colors.neutral[900], 0.9),
  },
}));

const Breadcrumbs = ({ entries }: BreadcrumbsProps) => {
  console.log("entries", entries);
  return (
    <StyledBreadcrumbs>
      {entries.map(({ url, label }) => (
        <Breadcrumb key={label} url={url} label={label} />
      ))}
    </StyledBreadcrumbs>
  );
};

export default Breadcrumbs;
