import React, { useState } from "react";
import { client, Row, Table } from "@clutch-sh/core";
import { Button, Container } from "@material-ui/core";

import { StartAbortExperiment, StartLatencyExperiment } from "./start-experiment";

function renderAbortData(experiment) {
  const ts = experiment.testSpecification;
  return (
    <Row
      key={experiment.id}
      data={[
        ts.abort.clusterPair.downstreamCluster,
        ts.abort.clusterPair.upstreamCluster,
        ts.abort.percent,
        ts.abort.httpStatus,
      ]}
    />
  );
}

function renderLatencyData(experiment) {
  const ts = experiment.testSpecification;
  return (
    <Row
      key={experiment.id}
      data={[
        ts.latency.clusterPair.downstreamCluster,
        ts.latency.clusterPair.upstreamCluster,
        ts.latency.percent,
        ts.latency.httpStatus,
      ]}
    />
  );
}

const ListExperiments = () => {
  const [experiments, setExperiments] = useState();

  if (!experiments) {
    client.post("/v1/experiments/get").then(response => {
      setExperiments(response?.data?.experiments || []);
    });

    return (
      <Table headings={["Downstream Cluster", "Upstream Cluster", "Percentage", "HTTP Status"]} />
    );
  }

  return (
    <Container>
      <Table
        data={experiments}
        headings={["Downstream Cluster", "Upstream Cluster", "Percentage", "HTTP Status"]}
      >
        {experiments.map(e => {
          if (e.testSpecification.abort) {
            return renderAbortData(e);
          }
          return renderLatencyData(e);
        })}
      </Table>
      <Button onClick={StartAbortExperiment}>Start Abort Experiment</Button>
      <Button onClick={StartLatencyExperiment}>Start Latency Experiment</Button>
    </Container>
  );
};

export default ListExperiments;
