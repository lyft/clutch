import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { ButtonGroup, client, Error, Row, Table } from "@clutch-sh/core";
import { Container } from "@material-ui/core";
import styled from "styled-components";

const ExperimentSpecificationData = ({ experiment, columns, mapping }) => {
  mapping = mapping || {};

  const converter = mapping[experiment.testConfig["@type"]];
  const converterExists = typeof converter !== "undefined";
  const model = converterExists ? converter(experiment.testConfig) : experiment;

  const data = columns.map(column => {
    if (column == "identifier") {
      return experiment.id;
    } else if (column in model) {
      return model[column];
    } else {
      return "Unknown";
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

const ListExperiments = ({ heading, columns, mapping, links }) => {
  const [experiments, setExperiments] = useState([]);
  const [error, setError] = useState("");

  const navigate = useNavigate();

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

  let columnNames = columns.map(name => name.toUpperCase());
  var links = links || []
  let buttons = links.map(link =>  {
    return {
      text: link.text,
      onClick: () => navigate(link.path)
    }
  })

  return (
    <Layout>
      {error && <Error message={error} />}
      <Table
        data={experiments}
        headings={columnNames}
      >
        {experiments.map(e => (
          <ExperimentSpecificationData key={e.id} experiment={e} columns={columns} mapping={mapping} />
        ))}
      </Table>
      <ButtonGroup
        buttons={buttons}
      />
    </Layout>
  );
};

export default ListExperiments;
