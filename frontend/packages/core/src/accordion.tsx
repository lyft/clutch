import * as React from "react";
import styled from "@emotion/styled";
import type { AccordionProps as MuiAccordionProps } from "@material-ui/core";
import {
  Accordion as MuiAccordion,
  AccordionActions as MuiAccordionActions,
  AccordionDetails as MuiAccordionDetails,
  AccordionSummary as MuiAccordionSummary,
  Divider as MuiDivider,
} from "@material-ui/core";
import AddIcon from "@material-ui/icons/Add";
import RemoveIcon from "@material-ui/icons/Remove";

const StyledAccordion = styled(MuiAccordion)({
  borderRadius: "4px",
  boxShadow: "none",
  border: "1px solid transparent",
  boxSizing: "content-box",

  "&.Mui-expanded": {
    margin: "0",
    boxShadow: "0px 4px 6px rgba(53, 72, 212, 0.2)",
    borderColor: "rgba(13, 16, 48, 0.12)",
  },

  ".MuiIconButton-root": {
    color: "rgba(13, 16, 48, 0.38)",
  },

  ".MuiIconButton-root.Mui-expanded": {
    color: "rgba(13, 16, 48, 0.6)",
  },

  ".MuiAccordionDetails-root": {
    padding: "16px",
    fontSize: "16px",
  },

  ".MuiAccordionActions-root": {
    padding: "0",
  },
});

const AccordionSummaryBase = ({ children, collapsible, expanded, ...props }) => {
  return (
    <MuiAccordionSummary
      expandIcon={collapsible ? expanded ? <RemoveIcon /> : <AddIcon /> : null}
      {...props}
    >
      {children}
    </MuiAccordionSummary>
  );
};

export const StyledAccordionSummary = styled(AccordionSummaryBase)({
  backgroundColor: "#fafafb",
  borderRadius: "4px",
  color: "#0d1030",
  height: "48px",

  "&:hover": {
    backgroundColor: "rgba(13, 16, 48, 0.03)",
  },

  "&:active": {
    backgroundColor: "rgba(13, 16, 48, 0.12)",
  },

  ".MuiAccordionSummary-content": {
    margin: "12px 0",
    fontSize: "16px",
  },

  "&.Mui-expanded": {
    backgroundColor: "#e1e4f9",
    minHeight: "48px",
  },
});

export interface AccordionProps extends Pick<MuiAccordionProps, "defaultExpanded"> {
  title?: string;
  collapsible?: boolean;
  children: React.ReactNode[];
}

export const Accordion = ({
  title,
  collapsible = true,
  defaultExpanded,
  children,
  ...props
}: AccordionProps) => {
  const [expanded, setExpanded] = React.useState(defaultExpanded);

  return (
    <StyledAccordion defaultExpanded={defaultExpanded} expanded={expanded} {...props}>
      <StyledAccordionSummary
        expanded={expanded}
        defaultExpanded={defaultExpanded}
        collapsible={collapsible}
        onClick={() => collapsible && setExpanded(v => !v)}
      >
        {title}
      </StyledAccordionSummary>
      {children}
    </StyledAccordion>
  );
};

export const AccordionActions = MuiAccordionActions;
export const AccordionDetails = MuiAccordionDetails;
export const AccordionDivider = MuiDivider;
