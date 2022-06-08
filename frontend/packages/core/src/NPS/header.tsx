import React from "react";
import {
  ClickAwayListener,
  Grow as MuiGrow,
  MenuList,
  Paper as MuiPaper,
  Popper as MuiPopper,
} from "@material-ui/core";
import ChatBubbleOutlineIcon from "@material-ui/icons/ChatBubbleOutline";
import { sortBy } from "lodash";

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
  marginLeft: "12px",
  zIndex: 1201,
});

const Paper = styled(MuiPaper)({
  width: "350px",
  boxShadow: "0px 15px 35px rgba(53, 72, 212, 0.2)",
  borderRadius: "8px",
});

const StyledFeedbackIcon = styled(IconButton)<{ $open: boolean }>(
  {
    color: "#ffffff",
    marginRight: "8px",
    padding: "12px",
    "&:hover": {
      background: "#2d3db4",
    },
    "&:active": {
      background: "#2938a5",
    },
  },
  props => ({
    background: props.$open ? "#2d3db4" : "unset",
  })
);

export const generateFeedbackTypes = (workflows: Workflow[]): SelectOption[] => {
  const feedbackTypes: SelectOption[] = [{ label: "General" }];

  const typeMap = {};

  workflows.forEach(workflow => {
    const { group, path, routes, displayName } = workflow;

    if (!typeMap[group]) {
      typeMap[group] = [];
    }

    typeMap[group].push(
      ...routes.map(route => ({
        label: route.displayName || displayName,
        value: `/${path}/${route.path}`.replace(/\/\/+/g, "/"),
      }))
    );
  });

  feedbackTypes.push(
    ...Object.keys(typeMap)
      .sort()
      .map(label => ({ label, group: sortBy(typeMap[label], ["label"]) }))
  );

  return feedbackTypes;
};

const HeaderFeedback = () => {
  const [open, setOpen] = React.useState<boolean>(false);
  const anchorRef = React.useRef(null);
  const { workflows } = useAppContext();

  const handleToggle = () => {
    setOpen(!open);
  };

  const handleClose = event => {
    // handler so that it wont close when selecting an item in the select
    if (event.target.localName === "body") {
      return;
    }
    if (anchorRef.current && anchorRef.current.contains(event.target)) {
      return;
    }
    setOpen(false);
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
                  <NPSFeedback origin="HEADER" feedbackTypes={generateFeedbackTypes(workflows)} />
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
