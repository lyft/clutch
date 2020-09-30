import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
import { ButtonGroup, client, Error } from "@clutch-sh/core";
import { Container } from "@material-ui/core";
import styled from "styled-components";

import { Column, ListView } from "./list-view";

const Layout = styled(Container)`
  padding: 5% 0;
`;

interface ExperimentTypeLinkProps {
  displayName: string;
  path: string;
}

interface ListExperimentsProps {
  columns: Column[];
  links: ExperimentTypeLinkProps[];
}

const ListExperiments: React.FC<ListExperimentsProps> = ({ columns, links }) => {
  const [experiments, setExperiments] = useState<
    IClutch.chaos.experimentation.v1.ListViewItem[] | undefined
  >(undefined);
  const [error, setError] = useState("");

  const navigate = useNavigate();

  const handleRowSelection = (event: any, item: IClutch.chaos.experimentation.v1.ListViewItem) => {
    navigate(`/experimentation/run/${item.identifier}`);
  };

  React.useEffect(() => {
    client
      .post("/v1/chaos/experimentation/getListView")
      .then(response => {
        setExperiments(response?.data?.items || []);
      })
      .catch(err => {
        setError(err.response.statusText);
      });
  }, []);

  const buttons = links.map(link => {
    return {
      text: link.displayName,
      onClick: () => navigate(link.path),
    };
  });

  return (
    <Layout>
      {error && <Error message={error} />}
      <ButtonGroup buttons={buttons} />
      <ListView
        columns={columns}
        items={experiments}
        onRowSelection={(event, item) => {
          handleRowSelection(event, item);
        }}
      />
    </Layout>
  );
};

export default ListExperiments;
