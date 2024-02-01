import * as React from "react";
import styled from "@emotion/styled";
import AddIcon from "@mui/icons-material/Add";
import RemoveIcon from "@mui/icons-material/Remove";
import type { AccordionProps as MuiAccordionProps, Theme } from "@mui/material";
import {
  Accordion as MuiAccordion,
  AccordionActions as MuiAccordionActions,
  AccordionDetails as MuiAccordionDetails,
  AccordionSummary as MuiAccordionSummary,
  alpha,
  Divider as MuiDivider,
  useControlled,
} from "@mui/material";

const StyledAccordion = styled(MuiAccordion)(({ theme }: { theme: Theme }) => ({
  borderRadius: "4px",
  boxShadow: "none",
  border: "1px solid transparent",
  boxSizing: "content-box",
  maxWidth: "100%",
  overflowWrap: "anywhere",

  "&.Mui-expanded": {
    boxShadow: `0px 4px 6px ${alpha(theme.palette.primary[600], 0.2)}`,
    borderColor: alpha(theme.palette.secondary[900], 0.12),
  },

  ".MuiIconButton-root": {
    color: alpha(theme.palette.secondary[900], 0.38),
  },

  ".MuiIconButton-root.Mui-expanded": {
    color: alpha(theme.palette.secondary[900], 0.6),
  },

  ".MuiAccordionDetails-root": {
    padding: "8px",
    fontSize: "16px",

    "> *": {
      margin: "8px",
    },
  },

  "&:before": {
    display: "none",
  },
}));

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

export const StyledAccordionSummary = styled(AccordionSummaryBase)(
  ({ theme }: { theme: Theme }) => ({
    backgroundColor: theme.palette.secondary[50],
    borderRadius: "4px",
    color: theme.palette.secondary[900],
    height: "48px",

    "&:hover": {
      backgroundColor: alpha(theme.palette.secondary[900], 0.03),
    },

    "&:active": {
      backgroundColor: alpha(theme.palette.secondary[900], 0.12),
    },

    ".MuiAccordionSummary-content": {
      margin: "12px 0",
      fontSize: "16px",
    },

    "&.Mui-expanded": {
      backgroundColor: theme.palette.primary[200],
      minHeight: "48px",
    },
  })
);

const StyledAccordionGroup = styled.div({
  width: "100%",

  ".MuiAccordion-root": {
    marginBottom: "16px",
  },

  ".MuiAccordion-root.Mui-expanded": {
    marginBottom: "16px",
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
    margin: "8px 8px",
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
