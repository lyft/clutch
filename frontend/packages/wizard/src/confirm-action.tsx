import React, { useState } from "react";
import {
  Button,
  ConfirmActionProps,
  Dialog,
  DialogActions,
  DialogContent,
  TextField,
  Typography,
} from "@clutch-sh/core";

interface ConfirmActionDialogProps extends ConfirmActionProps {
  open: boolean;
  onCancel: () => void;
}

const ConfirmAction: React.FC<ConfirmActionDialogProps> = ({
  open,
  title,
  description,
  confirmationText,
  onConfirm,
  onCancel,
}) => {
  const [input, setInput] = useState("");

  const handleConfirm = () => {
    if (input === confirmationText) {
      onConfirm();
    }
  };

  return (
    <Dialog title={title} open={open} onClose={onCancel}>
      <DialogContent>
        <Typography variant="body1">{description}</Typography>
        <TextField
          label="Confirmation"
          value={input}
          onChange={e => setInput(e.target.value)}
          fullWidth
        />
      </DialogContent>
      <DialogActions>
        <Button text="Cancel" onClick={onCancel} variant="secondary" />
        <Button
          text="Confirm"
          onClick={handleConfirm}
          variant="danger"
          disabled={input !== confirmationText}
        />
      </DialogActions>
    </Dialog>
  );
};

export default ConfirmAction;
