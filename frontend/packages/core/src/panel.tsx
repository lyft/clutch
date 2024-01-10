import React from "react";
import styled from "@emotion/styled";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import { Accordion as MuiExpansionPanel, AccordionDetails, AccordionSummary } from "@mui/material";

import { Typography } from "./typography";

const FullWidthExpansionPanel = styled(MuiExpansionPanel)`
  width: 100%;
`;

export interface AccordionProps {
  heading: string;
  summary: string;
  expanded?: boolean;
}

/** TODO: Combine with accordion */
const Accordion: React.FC<AccordionProps> = ({ heading, summary, expanded, children }) => {
  return (
    <FullWidthExpansionPanel defaultExpanded={expanded}>
      <AccordionSummary expandIcon={<ExpandMoreIcon />}>
        <Typography variant="body2">{heading}</Typography>
        <div style={{ flexGrow: 1 }} />
        <Typography variant="body2">{summary}</Typography>
      </AccordionSummary>
      <AccordionDetails>{children}</AccordionDetails>
    </FullWidthExpansionPanel>
  );
};

export default Accordion;
