import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { ButtonGroup, client, Error, Row, Table } from "@clutch-sh/core";
import { Container } from "@material-ui/core";
import styled from "styled-components";

const ExperimentSpecificationData = ({ experiment, columns, mapping }) => {
  const specification = experiment.testConfig;

  var data = [];  
  columns.forEach(column => {
    var item;
    if (column in mapping) {
      item = experiment.testConfig[mapping[column]]
    } else {
      item = experiment.testConfig[column]
    }

    if (typeof item === 'undefined') {
      data.push("Unknown")
    } else {
      data.push(item);
    }
  });

  return (
    <Row
      data={
        data
      }
    />
  );
};

const Layout = styled(Container)`
  padding: 5% 0;
`;

const ListExperiments = ({ heading, columns, mapping }) => {
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
      .post("/v1/experiments/get", { convert: true })
      .then(response => {
        setExperiments(response?.data?.experiments || []);
      })
      .catch(err => {
        setError(err.response.statusText);
      });
  }

  let column_names = columns.map(name => name.toUpperCase());

  return (
    <Layout>
      {error && <Error message={error} />}
      <Table
        data={experiments}
        headings={column_names}
      >
        {experiments.map(e => (
          <ExperimentSpecificationData key={e.id} experiment={e} columns={columns} mapping={mapping} />
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
