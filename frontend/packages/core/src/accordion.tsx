import * as React from "react";
import styled from "@emotion/styled";
import type { AccordionProps as MuiAccordionProps } from "@material-ui/core";
import {
  Accordion as MuiAccordion,
  AccordionActions as MuiAccordionActions,
  AccordionDetails as MuiAccordionDetails,
  AccordionSummary as MuiAccordionSummary,
  Divider as MuiDivider,
  useControlled,
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
    padding: "8px",
    fontSize: "16px",
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

const StyledAccordionGroup = styled.div({
  width: "100%",

  ".MuiAccordion-root": {
    marginBottom: "16px",
  },

  ".MuiAccordion-root.Mui-expanded": {
    marginBottom: "16px",
  },

  ".MuiAccordion-root:before": {
    display: "none",
  },
});

export interface AccordionProps extends Pick<MuiAccordionProps, "defaultExpanded" | "expanded"> {
  title?: string;
  collapsible?: boolean;
  children: React.ReactNode;
  onClick?: React.MouseEventHandler;
}

export const Accordion = ({
  title,
  collapsible = true,
  defaultExpanded,
  expanded: expandedProp,
  onClick: onClickProp,
  children,
  ...props
}: AccordionProps) => {
  const [expanded, setExpanded] = useControlled({
    controlled: expandedProp,
    default: defaultExpanded,
    name: "Clutch Accordion",
    state: "expanded",
  });

  const handleClick = (e: React.MouseEvent) => {
    if (collapsible) {
      setExpanded(!expanded);
    }
    if (onClickProp) {
      onClickProp(e);
    }
  };

  return (
    <StyledAccordion defaultExpanded={defaultExpanded} expanded={expanded} {...props}>
      <StyledAccordionSummary expanded={expanded} collapsible={collapsible} onClick={handleClick}>
        {title}
      </StyledAccordionSummary>
      {children}
    </StyledAccordion>
  );
};

export interface AccordionGroupProps {
  children?: React.ReactElement<AccordionProps> | React.ReactElement<AccordionProps>[];
  defaultExpandedIdx?: number;
}

export const AccordionGroup = ({ children, defaultExpandedIdx }: AccordionGroupProps) => {
  const [expandedIdx, setExpandedIdx] = React.useState(defaultExpandedIdx ?? -1);

  return (
    <StyledAccordionGroup>
      {
        // Clone each accordion as a controlled component.
        React.Children.map(children, (child, idx) =>
          React.cloneElement(child, {
            ...child.props,
            expanded: idx === expandedIdx,
            onClick: () => {
              setExpandedIdx(idx === expandedIdx ? -1 : idx);
            },
          })
        )
      }
    </StyledAccordionGroup>
  );
};

const StyledAccordionDetails = styled(MuiAccordionDetails)({
  display: "flex",
  flexWrap: "wrap",
  boxSizing: "border-box",
  "> *": {
    padding: "8px 8px",

  },
  ".MuiFormLabel-root": {
    padding: "inherit",
  },
});

const StyledAccordionActions = styled(MuiAccordionActions)({
  padding: "8px",
  "> *": {
    margin: "8px 8px",
  },
});

export const AccordionActions = StyledAccordionActions;
export const AccordionDetails = StyledAccordionDetails;
export const AccordionDivider = MuiDivider;
