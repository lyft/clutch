import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { ExpansionPanel, Row, Table } from "@clutch-sh/core";

interface ListenersProps {
  listeners: IClutch.envoytriage.v1.IListeners;
}

const Listeners: React.FC<ListenersProps> = ({ listeners }) => {
  const [statuses, setStatuses] = React.useState([]);
  const [summary, setSummary] = React.useState("");

  React.useEffect(() => {
    setStatuses(listeners.listenerStatuses);
  }, [listeners]);

  React.useEffect(() => {
    setSummary(`(${statuses.length} found)`);
  }, [statuses]);

  return (
    <ExpansionPanel heading="Listeners" summary={summary}>
      <Table headings={["Name", "Local Address"]}>
        {statuses.map(listener => (
          <Row key={listener.name} data={[listener.name, listener.localAddress]} />
        ))}
      </Table>
    </ExpansionPanel>
  );
};

export default Listeners;
