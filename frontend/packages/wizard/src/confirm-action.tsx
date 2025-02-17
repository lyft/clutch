import React from "react";
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  Typography,
  useWizardContext,
} from "@clutch-sh/core";

export interface ConfirmActionProps {
  title: string;
  description: React.ReactNode | string;
  actionLabel?: string;
  onConfirm?: () => void;
  onCancel?: () => void;
  submitDisabled?: boolean;
}

const ConfirmAction: React.FC<ConfirmActionProps> = ({
  title,
  description,
  actionLabel,
  onConfirm,
  onCancel,
  submitDisabled = false,
}) => {
  const { confirmActionOpen: open, setConfirmActionOpen: setOpen, onSubmit } = useWizardContext();

  const handleConfirm = () => {
    onSubmit();
    setOpen(false);
    if (onConfirm) {
      onConfirm();
    }
  };

  const handleCancel = () => {
    setOpen(false);
    if (onCancel) {
      onCancel();
    }
  };

  return (
    <Dialog title={title} open={open} onClose={handleCancel}>
      <DialogContent>
        {typeof description === "string" ? (
          <Typography variant="body1">{description}</Typography>
        ) : (
          description
        )}
      </DialogContent>
      <DialogActions>
        <Button text="Cancel" onClick={handleCancel} variant="neutral" />
        <Button
          text={actionLabel || "Confirm"}
          onClick={handleConfirm}
          variant="danger"
          disabled={submitDisabled}
        />
      </DialogActions>
    </Dialog>
  );
};

export default ConfirmAction;
