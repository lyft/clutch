import React from "react";
import {
  Accordion as MuiExpansionPanel,
  AccordionDetails,
  AccordionSummary,
  Typography,
} from "@material-ui/core";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import styled from "styled-components";

const FullWidthExpansionPanel = styled(MuiExpansionPanel)`
  width: 100%;
`;

export interface ExpansionPanelProps {
  heading: string;
  summary: string;
  expanded?: boolean;
}

const ExpansionPanel: React.FC<ExpansionPanelProps> = ({
  heading,
  summary,
  expanded,
  children,
}) => {
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

export default ExpansionPanel;
