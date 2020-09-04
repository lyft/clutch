import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
import { ButtonGroup, client, Error, Row, Table } from "@clutch-sh/core";
import { Container } from "@material-ui/core";
import styled from "styled-components";

interface ExperimentationDataProps {
  experiment: IClutch.chaos.experimentation.v1.Experiment;
  columns: string[];
  experimentTypes: any;
}

const ExperimentData: React.FC<ExperimentationDataProps> = ({
  experiment,
  columns,
  experimentTypes,
}) => {
  const types = experimentTypes || {};
  // Check for a configuration describing how a test of a given type should be displayed within the list view.
  if (!Object.prototype.hasOwnProperty.call(types, experiment.config["@type"])) {
    const data = columns.map(() => {
      return experiment.config["@type"];
    });
    return <Row data={data} />;
  }

  const registeredExperimentType = types[experiment.config["@type"]];

  const navigate = useNavigate();

  const mapperExists = Object.prototype.hasOwnProperty.call(registeredExperimentType, "mapping");
  const model = mapperExists ? registeredExperimentType.mapping(experiment.config) : experiment;

  const data = columns.map(column => {
    let value: string;
    if (column === "identifier") {
      value = experiment.id.toString();
    } else if (Object.prototype.hasOwnProperty.call(model, column)) {
      value = model[column];
    }

    return value ?? "Unknown";
  });

  return <Row hover onClick={(e) => {  navigate("/experimentation/view/"+experiment.id) }} data={data} />;
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
    const configuration = types[experimentType];
    if (Object.prototype.hasOwnProperty.call(configuration, "links")) {
      links = links.concat(configuration.links);
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
          <ExperimentData key={e.id} experiment={e} columns={columns} experimentTypes={types} />
        ))}
      </Table>
      <ButtonGroup buttons={buttons} />
    </Layout>
  );
};

export default ListExperiments;
