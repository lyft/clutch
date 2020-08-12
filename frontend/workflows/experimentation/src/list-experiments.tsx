import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
import { ButtonGroup, client, Error, Row, Table } from "@clutch-sh/core";
import { Container } from "@material-ui/core";
import styled from "styled-components";

interface ExperimentationSpecificationDataProps {
  experiment: IClutch.chaos.experimentation.v1.Experiment;
  columns: string[];
  experimentTypes: any;
}

const ExperimentSpecificationData: React.FC<ExperimentationSpecificationDataProps> = ({
  experiment,
  columns,
  experimentTypes,
}) => {
  const types = experimentTypes || {};
  const registeredExperimentType = types[experiment.testConfig["@type"]];
  const isExperimentTypeRegistered = typeof registeredExperimentType !== "undefined";
  if (!isExperimentTypeRegistered) {
    const data = columns.map(() => {
      return experiment.testConfig["@type"];
    });
    return <Row data={data} />;
  }

  const mapper = registeredExperimentType.mapping;
  const mapperExists = typeof mapper !== "undefined";
  const model = mapperExists ? mapper(experiment.testConfig) : experiment;

  const data = columns.map(column => {
    let value: string;
    if (column === "identifier") {
      value = experiment.id.toString();
    } else if (column in model) {
      value = model[column];
    }

    return value ?? "Unknown";
  });

  return <Row data={data} />;
};

const Layout = styled(Container)`
  padding: 5% 0;
`;

interface ExperimentTypeLinkProps {
  displayName: string;
  path: string;
}

interface ExperimentTypeProps {
  mapping: string;
  links: ExperimentTypeLinkProps[];
}

interface ListExperimentsProps {
  columns: [string];
  experimentTypes: ExperimentTypeProps[];
}

const ListExperiments: React.FC<ListExperimentsProps> = ({ columns, experimentTypes }) => {
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

  const types = experimentTypes || [];
  let links: ExperimentTypeLinkProps[] = [];
  Object.keys(types).forEach(experimentType => {
    const specification = types[experimentType];
    if (typeof specification.links !== "undefined") {
      links = links.concat(specification.links);
    }
  });

  const columnNames = columns.map(name => name.toUpperCase());
  const buttons = links.map(link => {
    return {
      text: link.displayName,
      onClick: () => navigate(link.path),
    };
  });

  return (
    <Layout>
      {error && <Error message={error} />}
      <Table headings={columnNames}>
        {experiments.map(e => (
          <ExperimentSpecificationData
            key={e.id}
            experiment={e}
            columns={columns}
            experimentTypes={types}
          />
        ))}
      </Table>
      <ButtonGroup buttons={buttons} />
    </Layout>
  );
};

export default ListExperiments;
