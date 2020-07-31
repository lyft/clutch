import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { ButtonGroup, client, Error, Row, Table } from "@clutch-sh/core";
import { Container } from "@material-ui/core";
import styled from "styled-components";

const ExperimentSpecificationData = ({ experiment }) => {
  const specification = experiment.testSpecification;
  const data = specification?.abort ? specification.abort : specification.latency;

  return (
    <Row
      data={[
        data.clusterPair.downstreamCluster,
        data.clusterPair.upstreamCluster,
        data.percent,
        data.httpStatus,
      ]}
    />
  );
};

const Layout = styled(Container)`
  padding: 5% 0;
`;

const ListExperiments = () => {
  const [experiments, setExperiments] = useState([]);
  const [error, setError] = useState("");

  const navigate = useNavigate();
  function handleClickStartAbortExperiment() {
    navigate("/experimentation/startabort");
  }

  function handleClickStartLatencyExperiment() {
    navigate("/experimentation/startlatency");
  }

  if (experiments.length === 0) {
    client
      .post("/v1/experiments/get")
      .then(response => {
        setExperiments(response?.data?.experiments || []);
      })
      .catch(err => {
        setError(err.response.statusText);
      });
  }

  return (
    <Layout>
      {error && <Error message={error} />}
      <Table
        data={experiments}
        headings={["Downstream Cluster", "Upstream Cluster", "Percentage", "HTTP Status"]}
      >
        {experiments.map(e => (
          <ExperimentSpecificationData key={e.id} experiment={e} />
        ))}
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
    </Layout>
  );
};

export default ListExperiments;
