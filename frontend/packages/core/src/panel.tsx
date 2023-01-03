import React from "react";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import {
  Accordion as MuiExpansionPanel,
  AccordionDetails,
  AccordionSummary,
  Typography,
} from "@mui/material";

import styled from "./styled";

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
        <Typography>{heading}</Typography>
        <div style={{ flexGrow: 1 }} />
        <Typography color="secondary">{summary}</Typography>
      </AccordionSummary>
      <AccordionDetails>{children}</AccordionDetails>
    </FullWidthExpansionPanel>
  );
};

export default Accordion;
