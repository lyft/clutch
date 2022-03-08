import React from "react";
import { useLocation } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
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
import { Toast } from "../Feedback";
import { TextField } from "../Input";
import { client } from "../Network";
import type { ClutchError } from "../Network/errors";
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
  width: "400px",
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
  const { workflows } = useAppContext();
  const {
    tempHydrateStore,
    data: { store },
  } = useStorageContext();
  const anchorRef = React.useRef(null);
  const location = useLocation();
  const [open, setOpen] = React.useState(false);
  const [shortLink, setShortLink] = React.useState<string | null>(null);
  const [validWorkflow, setValidWorkflow] = React.useState<boolean>(false);
  const [error, setError] = React.useState<ClutchError | null>(null);

  const checkValidWorkflow = () => {
    const workflow = workflowByRoute(workflows, location.pathname);

    setValidWorkflow(workflow?.shortLink ?? false);
  };

  // will trigger on a location change, emptying out our temporary storage and rechecking the workflow
  React.useEffect(() => {
    if (workflows.length) {
      checkValidWorkflow();
    }

    store(null, null, {});
  }, [location]);

  // Used for initial load to verify that once our workflows have loaded we are on a valid workflow
  React.useEffect(() => {
    if (workflows.length) {
      checkValidWorkflow();
    }
  }, [workflows]);

  const handleToggle = () => {
    setOpen(!open);
    setShortLink(null);
  };

  const handleClose = event => {
    if (anchorRef.current && anchorRef.current.contains(event.target)) {
      return;
    }
    setOpen(false);
  };

  // Will rotate our object into an array of type IShareableState to send in the API request
  const rotateStore = (): IClutch.shortlink.v1.IShareableState[] =>
    Object.keys(tempHydrateStore).map(key => ({ key, state: tempHydrateStore[key] }));

  const generateShortLink = () => {
    const requestData: IClutch.shortlink.v1.ICreateRequest = {
      path: `${location.pathname}${location.search}`,
      state: rotateStore(),
    };

    client
      .post("/v1/shortlink/create", requestData)
      .then(response => {
        const { hash } = response.data as IClutch.shortlink.v1.ICreateResponse;
        setShortLink(`${window.location.origin}/sl/${hash}`);
      })
      .catch((err: ClutchError) => {
        console.warn("failed to generate short link", err); // eslint-disable-line
        setError(err);
      });
  };

  if (!validWorkflow) {
    return null;
  }

  return (
    <>
      {error && (
        <Toast severity="error" onClose={() => setError(null)}>
          Unable to generate shortlink
        </Toast>
      )}
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
