import * as React from "react";
import styled from "@emotion/styled";
import type { DialogProps as MuiDialogProps } from "@material-ui/core";
import {
  Dialog as MuiDialog,
  DialogActions as MuiDialogActions,
  DialogContent as MuiDialogContent,
  DialogTitle as MuiDialogTitle,
  IconButton as MuiIconButton,
  Paper,
} from "@material-ui/core";
import CloseIcon from "@material-ui/icons/Close";

const DialogPaper = styled(Paper)({
  border: "1px solid rgba(13, 16, 48, 0.1)",
  boxShadow: "0px 10px 24px rgba(35, 48, 143, 0.3)",
  boxSizing: "border-box",
  backgroundColor: "#FFFFFF",
});

const DialogTitle = styled(MuiDialogTitle)({
  display: "flex",
  justifyContent: "space-between",
  fontSize: "20px",
  padding: "12px 12px 0 32px",
  fontWeight: 400,
  color: "#0D1030",
});

const DialogTitleText = styled.div({
  padding: "14px 0 0 0",
});

const IconButton = styled(MuiIconButton)({
  height: "12px",
  width: "12px",
  color: "#0D1030",
});

const DialogContent = styled(MuiDialogContent)({
  padding: "16px 32px 32px 32px",
  fontSize: "16px",
  fontWeight: 400,
  color: "rgba(13, 16, 48, 0.6)",
  "> *": {
    margin: "16px 0 0 0",
  },
});

const DialogActions = styled(MuiDialogActions)({
  borderTop: "1px solid rgba(13, 16, 48, 0.12)",
  padding: "0 8px",
  "> *": {
    margin: "16px 8px 16px 8px",
  },
});

export interface DialogContentProps {}
export interface DialogActionsProps {}
export type DialogCloseReasons = "closeButtonClick" | "escapeKeyDown" | "backdropClick";
export interface DialogProps extends Pick<MuiDialogProps, "open"> {
  title: string;
  children:
    | React.ReactElement<DialogContentProps>
    | React.ReactElement<DialogContentProps>[]
    | React.ReactElement<DialogActionsProps>
    | React.ReactElement<DialogActionsProps>[];
  onClose: (event?: object, reason?: DialogCloseReasons) => void;
}

const Dialog = ({ title, children, open, onClose }: DialogProps) => (
  <MuiDialog PaperComponent={DialogPaper} open={open} onClose={onClose}>
    <DialogTitle disableTypography>
      <DialogTitleText>{title}</DialogTitleText>
      <IconButton onClick={e => onClose(e, "closeButtonClick")}>
        <CloseIcon />
      </IconButton>
    </DialogTitle>
    {children}
  </MuiDialog>
);

export { Dialog, DialogActions, DialogContent };
