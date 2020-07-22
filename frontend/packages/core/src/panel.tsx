import React from "react";
import {
  ExpansionPanel as MuiExpansionPanel,
  ExpansionPanelDetails,
  ExpansionPanelSummary,
  Typography,
} from "@material-ui/core";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import styled from "styled-components";

const FullWidthExpansionPanel = styled(MuiExpansionPanel)`
  width: 100%;
`;

interface ExpansionPanelProps {
  heading: string;
  summary: string;
  expanded: boolean;
}

const ExpansionPanel: React.FC<ExpansionPanelProps> = ({
  heading,
  summary,
  expanded,
  children,
}) => {
  return (
    <FullWidthExpansionPanel expanded={expanded}>
      <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
        <Typography>{heading}</Typography>
        <div style={{ flexGrow: 1 }} />
        <Typography color="secondary">{summary}</Typography>
      </ExpansionPanelSummary>
      <ExpansionPanelDetails>{children}</ExpansionPanelDetails>
    </FullWidthExpansionPanel>
  );
};

export default ExpansionPanel;
