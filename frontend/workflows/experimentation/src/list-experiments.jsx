import React, { useState } from "react";
import { client, Row, Table } from "@clutch-sh/core";

const ListExperiments = () => {
  const [experiments, setExperiments] = useState();

  if (!experiments) {
    client.post("/v1/experiments/get").then(response => {
      setExperiments(response?.data?.experiments || []);
    });

    return <Table headings={["Downstream Cluster", "Upstream Cluster"]} />;
  }

  return (
    <Table
      data={experiments}
      headings={["Downstream Cluster", "Upstream Cluster", "Percentage", "HTTP Status"]}
    >
      {experiments.map(e => (
        <Row
          key={e.id}
          data={[
            e.testSpecification.abort.clusterPair.downstreamCluster,
            e.testSpecification.abort.clusterPair.upstreamCluster,
            e.testSpecification.abort.percent,
            e.testSpecification.abort.httpStatus,
          ]}
        />
      ))}
    </Table>
  );
};

export default ListExperiments;
