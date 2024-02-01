import React from "react";
import ChatBubbleOutlineIcon from "@mui/icons-material/ChatBubbleOutline";
import {
  alpha,
  ClickAwayListener,
  Grow as MuiGrow,
  MenuList,
  Paper as MuiPaper,
  Popper as MuiPopper,
  Theme,
} from "@mui/material";
import { get, sortBy } from "lodash";

import type { Workflow } from "../AppProvider/workflow";
import { IconButton } from "../button";
import { useAppContext } from "../Contexts";
import type { SelectOption } from "../Input";
import styled from "../styled";

import NPSFeedback from "./feedback";

const Grow = styled(MuiGrow)((props: { placement: string }) => ({
  transformOrigin: props.placement,
}));

const Popper = styled(MuiPopper)({
  padding: "0 12px",
  offset: "12px",
  zIndex: 1201,
});

const Paper = styled(MuiPaper)(({ theme }: { theme: Theme }) => ({
  width: "350px",
  boxShadow: `0px 15px 35px ${alpha(theme.palette.primary[600], 0.2)}`,
  borderRadius: "8px",
}));

const StyledFeedbackIcon = styled(IconButton)<{ $open: boolean }>(
  ({ theme }: { theme: Theme }) => ({
    color: theme.palette.contrastColor,
    marginRight: "8px",
    padding: "12px",
    "&:hover": {
      background: theme.palette.primary[600],
    },
    "&:active": {
      background: theme.palette.primary[700],
    },
  }),
  props => ({ theme }: { theme: Theme }) => ({
    background: props.$open ? theme.palette.primary[600] : "unset",
  })
);

export const generateFeedbackTypes = (workflows: Workflow[]): SelectOption[] => {
  const feedbackTypes: SelectOption[] = [{ label: "General" }];

  const typeMap = {};

  workflows.forEach(workflow => {
    const { group, path, routes, displayName } = workflow;

    routes.forEach(route => {
      const additionalNPS = get(route, "componentProps.additionalNPS", []);

      const showRoute = route.hideNav === undefined || route.hideNav === false;
      if (showRoute || additionalNPS.length > 0) {
        if (!typeMap[group]) {
          typeMap[group] = [];
        }

        if (showRoute) {
          typeMap[group].push({
            label: route.displayName || displayName,
            value: `/${path}/${route.path}`.replace(/\/\/+/g, "/"),
          });
        }

        if (additionalNPS.length > 0) {
          typeMap[group].push(...additionalNPS);
        }
      }
    });
  });

  feedbackTypes.push(
    ...Object.keys(typeMap)
      .sort()
      .map(label => ({ label, group: sortBy(typeMap[label], ["label"]) }))
  );

  return feedbackTypes;
};

/**
 * An NPS Header component which will render an icon in the banner and when clicked
 * will ask the user to provide feedback
 */
const HeaderFeedback = () => {
  const [open, setOpen] = React.useState<boolean>(false);
  const anchorRef = React.useRef(null);
  const { workflows, triggerHeaderItem, triggeredHeaderData } = useAppContext();
  const [defaultFeedbackOption, setDefaultFeedbackOption] = React.useState<string>();
  const timer = React.useRef(null);

  const handleToggle = () => {
    setOpen(!open);
  };

  const timedClose = () => {
    timer.current = setTimeout(() => {
      setOpen(false);
    }, 1500);
  };

  React.useEffect(() => {
    if (triggeredHeaderData && triggeredHeaderData.NPS) {
      setDefaultFeedbackOption((triggeredHeaderData.NPS.defaultFeedbackOption as string) ?? "");
      setOpen(true);
    }
  }, [triggeredHeaderData]);

  const handleClose = event => {
    // handler so that it wont close when selecting an item in the select
    if (event.target.localName === "body") {
      return;
    }
    if (anchorRef.current && anchorRef.current.contains(event.target)) {
      return;
    }
    // handler for the NPS Banner button so that it doesn't reset the headerLink
    if (event.target.id !== "nps-banner-button") {
      triggerHeaderItem && triggerHeaderItem("NPS", undefined);
      setOpen(false);
      clearTimeout(timer.current);
    }
  };

  return (
    <>
      <StyledFeedbackIcon
        variant="neutral"
        ref={anchorRef}
        aria-controls={open ? "header-feedback" : undefined}
        $open={open}
        aria-haspopup="true"
        onClick={handleToggle}
        edge="end"
        id="headerFeedbackIcon"
      >
        <ChatBubbleOutlineIcon />
      </StyledFeedbackIcon>
      <Popper open={open} anchorEl={anchorRef.current} transition placement="bottom-end">
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            placement={placement === "bottom" ? "center top" : "center bottom"}
          >
            <Paper>
              <ClickAwayListener onClickAway={handleClose}>
                <MenuList autoFocusItem={open} id="options">
                  <NPSFeedback
                    origin="HEADER"
                    feedbackTypes={generateFeedbackTypes(workflows)}
                    defaultFeedbackOption={defaultFeedbackOption}
                    onSubmit={timedClose}
                  />
                </MenuList>
              </ClickAwayListener>
            </Paper>
          </Grow>
        )}
      </Popper>
    </>
  );
};

export default HeaderFeedback;
