import React from "react";
import { useLocation } from "react-router-dom";
import {
  ClickAwayListener,
  Grow as MuiGrow,
  MenuList,
  Paper as MuiPaper,
  Popper as MuiPopper,
} from "@material-ui/core";
import FileCopyIcon from "@material-ui/icons/FileCopy";
import LinkIcon from "@material-ui/icons/Link";

import { IconButton } from "../button";
import { useStorageContext } from "../Contexts";
import { TextField } from "../Input";
import styled from "../styled";

const Grow = styled(MuiGrow)((props: { placement: string }) => ({
  transformOrigin: props.placement,
}));

const Popper = styled(MuiPopper)({
  padding: "0 12px",
  marginLeft: "12px",
  zIndex: 1201,
});

const Paper = styled(MuiPaper)({
  width: "450px",
  height: "100px",
  padding: "15px",
  boxShadow: "0px 15px 35px rgba(53, 72, 212, 0.2)",
  borderRadius: "8px",
});

const StyledLinkIcon = styled(IconButton)<{ $open: boolean }>(
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

const ShortLinker = () => {
  const [open, setOpen] = React.useState(false);
  const anchorRef = React.useRef(null);
  const location = useLocation();
  const {
    tempHydrateStore,
    data: { store },
  } = useStorageContext();
  const [shortLink, setShortLink] = React.useState("https://clutch.lyft.net/sl/1234");

  React.useEffect(() => {
    console.log("CLEARING STORAGE");
    store(null, null, {});
  }, [location]);

  React.useEffect(() => {
    console.log("CHECKING STORAGE", tempHydrateStore);
  }, [tempHydrateStore]);

  // on click, take current route and temp hydrate data and do stuff

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
      <StyledLinkIcon
        variant="neutral"
        ref={anchorRef}
        aria-controls={open ? "header-feedback" : undefined}
        $open={open}
        aria-haspopup="true"
        onClick={handleToggle}
        edge="end"
        id="headerFeedbackIcon"
      >
        <LinkIcon />
      </StyledLinkIcon>
      <Popper open={open} anchorEl={anchorRef.current} transition placement="bottom-end">
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            placement={placement === "bottom" ? "center top" : "center bottom"}
          >
            {/* Get an error about a failed prop type without the fragment, needs further investigation */}
            <>
              <Paper>
                <ClickAwayListener onClickAway={handleClose}>
                  <MenuList autoFocusItem={open} id="options">
                    <TextField
                      disabled
                      readOnly
                      value={shortLink}
                      endAdornment={<FileCopyIcon />}
                    />
                  </MenuList>
                </ClickAwayListener>
              </Paper>
            </>
          </Grow>
        )}
      </Popper>
    </>
  );
};

export default ShortLinker;
