import React from "react";

import { Link } from "../link";
import styled from "../styled";
import { Typography } from "../typography";

import { BreadcrumbEntry } from "./types";

const StyledTypography = styled(Typography)({
  fontWeight: 500,
});

const Breadcrumb = ({ label, url }: BreadcrumbEntry) => {
  return url ? (
    <Link href={url} target="_self" whiteSpace="nowrap">
      <StyledTypography variant="caption2" color="inherit">
        {label}
      </StyledTypography>
    </Link>
  ) : (
    <StyledTypography variant="caption2">{label}</StyledTypography>
  );
};

export default Breadcrumb;
