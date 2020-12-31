import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Table, TableRow } from "@clutch-sh/core";

interface ListenersProps {
  listeners: IClutch.envoytriage.v1.IListeners;
}

const Listeners: React.FC<ListenersProps> = ({ listeners }) => {
  const [statuses, setStatuses] = React.useState([]);

  React.useEffect(() => {
    setStatuses(listeners.listenerStatuses);
  }, [listeners]);

  return (
    <Table headings={["Name", "Local Address"]}>
      {statuses.map(listener => (
        <TableRow key={listener.name}>
          {listener.name}
          {listener.localAddress}
        </TableRow>
      ))}
    </Table>
  );
};

export default Listeners;
