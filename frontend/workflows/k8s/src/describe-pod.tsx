import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { MetadataTable, Resolver, Table, TableRow, useWizardContext } from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import _ from "lodash";

import type { ResolverChild, WorkflowProps } from ".";

const PodIdentifier: React.FC<ResolverChild> = ({ resolverType }) => {
  const { onSubmit } = useWizardContext();
  const resolvedResourceData = useDataLayout("resourceData");
  const resolverInput = useDataLayout("resolverInput");

  const onResolve = ({ results, input }) => {
    // Decide how to process results.
    resolvedResourceData.assign(results[0]);
    resolverInput.assign(input);
    onSubmit();
  };

  return <Resolver type={resolverType} searchLimit={1} onResolve={onResolve} />;
};

const PodDetails: React.FC<WizardChild> = () => {
  const resourceData = useDataLayout("resourceData");
  const instance = resourceData.displayValue() as IClutch.k8s.v1.Pod;
  const { containers } = instance;

  return (
    <WizardStep error={resourceData.error} isLoading={resourceData.isLoading}>
      <strong>Pod Details</strong>
      <MetadataTable
        data={[
          { name: "Name", value: instance.name },
          { name: "Cluster", value: instance.cluster },
          { name: "Namespace", value: instance.namespace },
          { name: "State", value: _.capitalize(instance.state.toString()) },
          { name: "Node IP Address", value: instance.nodeIp },
          { name: "Pod IP Address", value: instance.podIp },
          {
            name: "Containers",
            value: (
              <Table stickyHeader headings={["Name", "State", "Restart Count"]}>
                {_.sortBy(containers, [
                  o => {
                    return o.name;
                  },
                ]).map(container => (
                  <TableRow key={container.name} defaultCellValue="nil">
                    {container.name}
                    {container.state}
                    {container.restartCount}
                  </TableRow>
                ))}
              </Table>
            ),
          },
        ]}
      />
    </WizardStep>
  );
};

const DescribePod: React.FC<WorkflowProps> = ({ heading, resolverType }) => {
  const dataLayout = {
    resolverInput: {},
    resourceData: {},
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <PodIdentifier name="Lookup" resolverType={resolverType} />
      <PodDetails name="Details" />
    </Wizard>
  );
};

export default DescribePod;
