import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
import { BaseWorkflowProps, ButtonGroup, client } from "@clutch-sh/core";

import PageLayout from "./core/page-layout";
import { Column, ListView } from "./list-view";

interface ExperimentTypeLinkProps {
  displayName: string;
  path: string;
}

interface ListExperimentsProps extends BaseWorkflowProps {
  columns: Column[];
  links: ExperimentTypeLinkProps[];
}

const ListExperiments: React.FC<ListExperimentsProps> = ({ heading, columns, links }) => {
  const [experiments, setExperiments] = useState<
    IClutch.chaos.experimentation.v1.ListViewItem[] | undefined
  >(undefined);
  const [error, setError] = useState("");

  const navigate = useNavigate();

  const handleRowSelection = (event: any, item: IClutch.chaos.experimentation.v1.ListViewItem) => {
    navigate(`/experimentation/run/${item.id}`);
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
    <PageLayout heading={heading} error={error}>
      <ButtonGroup buttons={buttons} />
      <ListView
        columns={columns}
        items={experiments}
        onRowSelection={(event, item) => {
          handleRowSelection(event, item);
        }}
      />
    </PageLayout>
  );
};

export default ListExperiments;
