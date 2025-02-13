import React from "react";
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
  actionLabel,
  onConfirm,
  onCancel,
}) => {
  const [input, setInput] = React.useState("");

  React.useEffect(() => {
    if (open) {
      setInput("");
    }
  }, [open]);

  const handleConfirm = () => {
    if (input === confirmationText) {
      onConfirm();
    }
  };

  return (
    <Dialog title={title} open={open} onClose={onCancel}>
      <DialogContent>
        {typeof description === "string" ? (
          <Typography variant="body1">{description}</Typography>
        ) : (
          description
        )}
        <TextField
          label="Confirmation"
          value={input}
          onChange={e => setInput(e.target.value)}
          fullWidth
        />
      </DialogContent>
      <DialogActions>
        <Button text="Cancel" onClick={onCancel} variant="neutral" />
        <Button
          text={actionLabel || "Confirm"}
          onClick={handleConfirm}
          variant="danger"
          disabled={input !== confirmationText}
        />
      </DialogActions>
    </Dialog>
  );
};

export default ConfirmAction;
