import React, { useState } from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import type { ClutchError } from "@clutch-sh/core";
import { BaseWorkflowProps, Button, ButtonGroup, client, useNavigate } from "@clutch-sh/core";

import PageLayout from "./core/page-layout";
import type { Column } from "./list-view";
import ListView from "./list-view";

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
  const [error, setError] = useState<ClutchError | undefined>(undefined);

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
      .catch((err: ClutchError) => {
        setError(err);
      });
  }, []);

  const buttons = links.map(link => (
    <Button text={link.displayName} key={link.path} onClick={() => navigate(link.path)} />
  ));

  return (
    <PageLayout heading={heading} error={error}>
      <ButtonGroup justify="center" border="bottom">
        {buttons}
      </ButtonGroup>
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
