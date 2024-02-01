import React from "react";
import { useLocation } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
import LinkIcon from "@mui/icons-material/Link";
import {
  alpha,
  ClickAwayListener,
  Grid,
  Grow as MuiGrow,
  MenuList,
  Paper as MuiPaper,
  Popper as MuiPopper,
  Theme,
} from "@mui/material";

import { generateShortLinkRoute } from "../AppProvider/short-link-proxy";
import { Button, ClipboardButton, IconButton } from "../button";
import { useAppContext, useShortLinkContext } from "../Contexts";
import type { HydratedData } from "../Contexts/workflow-storage-context/types";
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

const Paper = styled(MuiPaper)(({ theme }: { theme: Theme }) => ({
  width: "400px",
  height: "100px",
  padding: "15px",
  boxShadow: `0px 5px 35px ${alpha(theme.palette.primary[400], 0.2)}`,
  borderRadius: "8px",
}));

const StyledLinkIcon = styled(IconButton)<{ $open: boolean }>(
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
  props => ({
    background: props.$open ? props.theme.palette.primary[600] : "unset",
  })
);

/**
 * Component that will display a Button to generate a short link
 * - Will only be displayed if the given workflow has defined the `shortLink` property
 * On click, the component will open a popper window with a Generate ShortLink button
 * On clicking the Generate ShortLink button, the component will do the following:
 * - Read in the data from the temporary storage in the ShortLinkContext
 * - Rotate it into a readable format for the API
 * - Send an API request asking for it to be stored
 * - If successful,
 * -     Will switch to a readable only input field with a clipboard button
 *        to copy the short link with generated hash
 *       (this is the only time it is ever displayed)
 * - If not successful,
 * -     Will display a toast error message for the user
 */
const ShortLinker = () => {
  const { workflows } = useAppContext();
  const { removeWorkflowSession, retrieveWorkflowSession } = useShortLinkContext();
  const [open, setOpen] = React.useState(false);
  const [shortLink, setShortLink] = React.useState<string | null>(null);
  const [validWorkflow, setValidWorkflow] = React.useState<boolean>(false);
  const [error, setError] = React.useState<ClutchError | null>(null);
  const anchorRef = React.useRef(null);
  const location = useLocation();

  /**
   * Will fetch the the workflow based on the current route
   * Then checks to see if the shortLink property of the workflow has been set.
   *  - If it has, we will set the validWorkflow state which will allow the component to render
   */
  const checkValidWorkflow = () => {
    const workflow = workflowByRoute(workflows, location.pathname);

    setValidWorkflow(workflow?.shortLink ?? false);
  };

  // will trigger on a location change, emptying out our temporary storage and rechecking the workflow
  React.useEffect(() => {
    if (workflows.length) {
      checkValidWorkflow();
    }

    removeWorkflowSession();
  }, [location]);

  // Used for initial load to verify that once our workflows have loaded we are on a valid workflow
  React.useEffect(() => {
    if (workflows.length) {
      checkValidWorkflow();
    }
  }, [workflows]);

  const handleToggle = () => {
    setOpen(o => !o);
    setShortLink(null);
  };

  const handleClose = event => {
    if (anchorRef.current && anchorRef.current.contains(event.target)) {
      return;
    }
    setOpen(false);
  };

  // Will rotate our object into an array of type IShareableState to send in the API request
  const rotateStore = (sessionStore: HydratedData): IClutch.shortlink.v1.IShareableState[] =>
    Object.keys(sessionStore).map(key => ({ key, state: sessionStore[key] }));

  const generateShortLink = () => {
    if (error) {
      setError(null);
    }

    const requestData: IClutch.shortlink.v1.ICreateRequest = {
      path: `${location.pathname}${location.search}`,
      state: rotateStore(retrieveWorkflowSession()),
    };

    client
      .post("/v1/shortlink/create", requestData)
      .then(response => {
        const { hash } = response.data as IClutch.shortlink.v1.ICreateResponse;
        setShortLink(generateShortLinkRoute(window.location.origin, hash));
      })
      .catch((err: ClutchError) => {
        setError(err);
        setOpen(false);
      });
  };

  if (!validWorkflow) {
    return null;
  }

  return (
    <>
      {error && (
        <Toast title="Generating Short Link" severity="error" onClose={() => setError(null)}>
          {error?.message}
        </Toast>
      )}
      <StyledLinkIcon
        variant="neutral"
        ref={anchorRef}
        aria-controls={open ? "header-shortlink" : undefined}
        $open={open}
        aria-haspopup="true"
        onClick={handleToggle}
        edge="end"
        id="headerShortLinkIcon"
      >
        <LinkIcon />
      </StyledLinkIcon>
      <Popper open={open} anchorEl={anchorRef.current} transition placement="bottom-end">
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            placement={placement === "bottom" ? "center top" : "center bottom"}
          >
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
          </Grow>
        )}
      </Popper>
    </>
  );
};

export default ShortLinker;
