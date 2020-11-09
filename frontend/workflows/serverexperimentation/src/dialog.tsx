import React from "react";
import Button from "@material-ui/core/Button";
import MuiDialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogTitle from "@material-ui/core/DialogTitle";
import Typography from "@material-ui/core/Typography";

interface ButtonProps {
  label: string;
  onAction: () => void;
}

interface DialogProps {
  title: string;
  content: string;
  open: boolean;
  onClose: () => void;
  buttons: ButtonProps[];
}

const Dialog: React.FC<DialogProps> = ({ title, content, open, onClose, buttons }) => {
  return (
    <MuiDialog open={open} onClose={onClose}>
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        <Typography gutterBottom>{content}</Typography>
      </DialogContent>
      <DialogActions>
        {buttons.map(buttonProps => {
          return (
            <Button key={buttonProps.label} onClick={buttonProps.onAction}>
              {buttonProps.label}
            </Button>
          );
        })}
      </DialogActions>
    </MuiDialog>
  );
};

export default Dialog;
