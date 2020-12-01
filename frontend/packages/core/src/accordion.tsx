import * as React from "react";
import styled from "@emotion/styled";
import { Accordion as MuiAccordion, AccordionDetails, AccordionActions, AccordionSummary as MuiAccordionSummary } from "@material-ui/core";

import AddCircleIcon from "@material-ui/icons/AddCircle";
import RemoveCircleIcon from "@material-ui/icons/RemoveCircle";


export const Accordion = styled(MuiAccordion)({
    borderRadius: "4px",
    boxShadow: "none",

    "&.Mui-expanded": {
        margin: "0",
        boxShadow: "0px 4px 6px rgba(53, 72, 212, 0.2)"
    },

    ".MuiIconButton-root": {
        color: "rgba(13, 16, 48, 0.38)",
    },

    ".MuiIconButton-root.Mui-expanded": {
        color: "rgba(13, 16, 48, 0.6)"
    },

    ".MuiAccordionDetails-root": {
        padding: "16px",
    },

    ".MuiAccordionActions-root": {
        padding: "16px",
    },

    ".MuiCollapse-wrapper": {
        borderRadius: "0 0 4px 4px",
        border: "1px solid #e1e4f9",
    }
});

const AccordionSummaryBase = ({ children, ...props }) => {
    const [expanded, setExpanded] = React.useState(false);
  
    const onClick = () => {
      setExpanded(!expanded);
    };
  
    return (
      <MuiAccordionSummary
        {...props}
        expandIcon={expanded ? <AddCircleIcon /> : <RemoveCircleIcon />}
        onClick={onClick}
      >
        {children}
      </MuiAccordionSummary>
    );
  };

export const AccordionSummary = styled(AccordionSummaryBase)({
    backgroundColor: "#fafafb",
    borderRadius: "4px 4px 0 0",

    ".MuiAccordionSummary-content": {
        margin: "12px 0",
        fontSize: "16px",
    },

    "&.Mui-expanded": {
        backgroundColor: "#e1e4f9",
        minHeight: "48px",
    }
})
