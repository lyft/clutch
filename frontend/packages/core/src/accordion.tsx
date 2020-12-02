import * as React from "react";
import styled from "@emotion/styled";
import type { AccordionProps as MuiAccordionProps } from "@material-ui/core";
import {
  Accordion as MuiAccordion,
  AccordionActions,
  AccordionDetails,
  AccordionSummary as MuiAccordionSummary,
  Divider,
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

export interface AccordionProps extends Pick<MuiAccordionProps, "defaultExpanded"> {
  title?: string;
  collapsible?: boolean;
  children: React.ReactNode[];
}

export const Accordion = ({
  title,
  collapsible,
  defaultExpanded,
  children,
  ...props
}: AccordionProps) => (
  <StyledAccordion defaultExpanded={defaultExpanded} {...props}>
    <StyledAccordionSummary defaultExpanded={defaultExpanded} collapsible={collapsible}>
      {title}
    </StyledAccordionSummary>
    {children}
  </StyledAccordion>
);

const AccordionSummaryBase = ({ children, collapsible, defaultExpanded, ...props }) => {
  const [expanded, setExpanded] = React.useState(defaultExpanded);

  return (
    <MuiAccordionSummary
      {...props}
      expandIcon={collapsible ? expanded ? <RemoveIcon /> : <AddIcon /> : null}
      onClick={() => setExpanded((v: boolean) => !v)}
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

export { AccordionActions, AccordionDetails, Divider as AccordionDivider };
