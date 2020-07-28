import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { ButtonGroup, client, Row, Table } from "@clutch-sh/core";
import { Container } from "@material-ui/core";

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
        ts.latency.durationMs,
      ]}
    />
  );
}

const ListExperiments = () => {
  const [experiments, setExperiments] = useState();

  const navigate = useNavigate();
  function handleClickStartAbortExperiment() {
    navigate("/experimentation/startabort");
  }

  function handleClickStartLatencyExperiment() {
    navigate("/experimentation/startlatency");
  }

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
      <ButtonGroup
        buttons={[
          {
            text: "Start Abort Experiment",
            onClick: () => handleClickStartAbortExperiment(),
          },
          {
            text: "Start Latency Experiment",
            onClick: () => handleClickStartLatencyExperiment(),
          },
        ]}
      />
    </Container>
  );
};

export default ListExperiments;
