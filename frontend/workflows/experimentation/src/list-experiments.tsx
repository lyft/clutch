import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
import type { BaseWorkflowProps } from "@clutch-sh/core";
import { ButtonGroup, client, Error, Row, Table } from "@clutch-sh/core";
import { Container } from "@material-ui/core";
import styled from "styled-components";

interface ExperimentationSpecificationDataProps {
  experiment: IClutch.chaos.experimentation.v1.Experiment;
}

const ExperimentSpecificationData: React.FC<ExperimentationSpecificationDataProps> = ({
  experiment,
}) => {
  const specification = experiment.testSpecification;
  const data = specification?.abort ? specification.abort : specification.latency;

  const defaultData = [
    data.clusterPair.downstreamCluster,
    data.clusterPair.upstreamCluster,
    data.percent,
  ];
  if (specification?.abort) {
    defaultData.push((data as IClutch.chaos.experimentation.v1.AbortFault).httpStatus);
  }

  return <Row data={defaultData} />;
};

const Layout = styled(Container)`
  padding: 5% 0;
`;

const ListExperiments: React.FC<BaseWorkflowProps> = () => {
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
      <Table headings={["Downstream Cluster", "Upstream Cluster", "Percentage", "HTTP Status"]}>
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
