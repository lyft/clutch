import React from "react";

import { Link } from "../link";
import styled from "../styled";
import { Typography } from "../typography";

const StyledLink = styled(Link)({
  whiteSpace: "nowrap",
});

const StyledTypography = styled(Typography)({
  fontWeight: 500,
});

const Breadcrumb = ({ label, url }) => {
  return url ? (
    <StyledLink href={url} target="_self">
      <StyledTypography variant="caption2" color="inherit">
        {label}
      </StyledTypography>
    </StyledLink>
  ) : (
    <StyledTypography variant="caption2">{label}</StyledTypography>
  );
};

export default Breadcrumb;
