import React from "react";
import { useLocation } from "react-router-dom";
import {
  ClickAwayListener,
  Grid,
  Grow as MuiGrow,
  MenuList,
  Paper as MuiPaper,
  Popper as MuiPopper,
} from "@material-ui/core";
import LinkIcon from "@material-ui/icons/Link";

import { Button, ClipboardButton, IconButton } from "../button";
import { useAppContext, useStorageContext } from "../Contexts";
import { TextField } from "../Input";
import { client } from "../Network";
import styled from "../styled";

import { workflowByRoute } from "./utils";

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
  const { workflows } = useAppContext();
  const {
    tempHydrateStore,
    data: { store },
  } = useStorageContext();
  const location = useLocation();
  const [shortLink, setShortLink] = React.useState<string | null>(null);

  React.useEffect(() => {
    if (workflows.length) {
    }
    // Will clear our temp storage on location change
    store(null, null, {});
  }, [location]);

  React.useEffect(() => {
    if (workflows.length) {
      console.log(workflowByRoute(workflows, location.pathname));
    }
  }, [workflows]);

  // on click, take current route and temp hydrate data and do stuff

  const handleToggle = () => {
    setOpen(!open);
    setShortLink(null);
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

  const generateShortLink = () => {
    // call api with tempHydrateStore and route
    // client
    //   .post("/v1/...", { route: window.location.href, data: tempHydrateStore })
    //   .then(response => {
    //     setShortLink(`${window.location.origin}/sl/${response.code}`);
    //   })
    //   .catch((error: ClutchError) => {
    //     console.warn("failed to generate short link", error); // eslint-disable-line
    //     // throw a toast???
    //   });
    // get code
    // set short link
    const code = 1234;
    setShortLink(`${window.location.origin}/sl/${code}`);
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
            <>
              <Paper>
                <ClickAwayListener onClickAway={handleClose}>
                  <MenuList autoFocusItem={open} id="options">
                    <Grid style={{ display: "flex", justifyContent: "center" }}>
                      {!shortLink && (
                        <Button onClick={generateShortLink} text="Generate Short Link" />
                      )}
                      {shortLink && (
                        <>
                          <TextField disabled readOnly value={shortLink} />
                          <ClipboardButton text={shortLink} tooltip="Copy Short Link" />
                        </>
                      )}
                    </Grid>
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
