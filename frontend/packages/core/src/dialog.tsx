import React from "react";
import MuiDialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogTitle from "@material-ui/core/DialogTitle";
import Typography from "@material-ui/core/Typography";

export interface DialogProps {
  title: string;
  content: string;
  open: boolean;
  onClose: () => void;
}

const Dialog: React.FC<DialogProps> = ({ title, content, open, onClose, children }) => {
  return (
    <MuiDialog open={open} onClose={onClose}>
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        <Typography gutterBottom>{content}</Typography>
      </DialogContent>
      <DialogActions>{children}</DialogActions>
    </MuiDialog>
  );
};

export default Dialog;
