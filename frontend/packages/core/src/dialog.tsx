import * as React from "react";
import styled from "@emotion/styled";
import type { DialogProps as MuiDialogProps } from "@material-ui/core";
import {
  Dialog as MuiDialog,
  DialogActions as MuiDialogActions,
  DialogContent as MuiDialogContent,
  DialogTitle as MuiDialogTitle,
  IconButton as MuiIconButton,
} from "@material-ui/core";
import CloseIcon from "@material-ui/icons/Close";

const DialogTitle = styled(MuiDialogTitle)({
  display: "flex",
  justifyContent: "space-between",
  fontSize: "16px",
  padding: "16px 16px 0 16px",
  fontWeight: 400,
  color: "#000000",
});

const IconButton = styled(MuiIconButton)({
  height: "16px",
  width: "16px",
  padding: "0",
  color: "#000000",
});

const DialogContent = styled(MuiDialogContent)({
  padding: "8px 16px",
  fontSize: "12px",
  fontWeight: 400,
  color: "rgba(13, 16, 48, 0.6)",
});

const DialogActions = styled(MuiDialogActions)({
  borderTop: "1px solid rgba(13, 16, 48, 0.12)",
  padding: "0 8px",
});

export interface DialogProps extends Pick<MuiDialogProps, "open"> {
  title: string;
  children: React.ReactNode;
  actions: React.ReactNode;
  onClose: () => void;
}

export const Dialog = ({ title, children, open, onClose, actions }: DialogProps) => (
  <MuiDialog open={open} onClose={onClose}>
    <DialogTitle disableTypography>
      <div>{title}</div>
      <IconButton onClick={onClose}>
        <CloseIcon fontSize="small" />
      </IconButton>
    </DialogTitle>
    <DialogContent>{children}</DialogContent>
    <DialogActions>{actions}</DialogActions>
  </MuiDialog>
);

export default Dialog;
