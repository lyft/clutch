import React from "react";

import { Dialog, DialogContent } from "../../dialog";
import type { ClutchError } from "../../Network/errors";
import Code from "../../text";

interface ErrorDetailsDialogProps {
  error: ClutchError;
  open: boolean;
  onClose: () => void;
}

const ErrorDetailsDialog = ({ error, open, onClose }: ErrorDetailsDialogProps) => {
  const prettyPrintError = JSON.stringify(error, undefined, 2);

  return (
    <Dialog title="Full Error Details" onClose={onClose} open={open}>
      <DialogContent>
        {/* TODO: Plumb request ID through once we have it */}
        <Code>{prettyPrintError}</Code>
      </DialogContent>
    </Dialog>
  );
};

export default ErrorDetailsDialog;
