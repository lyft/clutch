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
        <div style={{ display: "flex", alignItems: "center" }}>
          <div>Request ID:&nbsp;</div>
          {/* TODO: Plumb request ID through */}
          <Code isCopiable={false}>ade67hea-4dbb-89hj-33ddbb87cae</Code>
        </div>
        <Code>{prettyPrintError}</Code>
      </DialogContent>
    </Dialog>
  );
};

export default ErrorDetailsDialog;
