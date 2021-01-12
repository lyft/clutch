import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Grid, MetadataTable, Paper } from "@clutch-sh/core";
import styled from "@emotion/styled";
import { Pie } from "@nivo/pie";

const SummaryCardTitle = styled.div({
  fontWeight: 600,
  fontSize: "14px",
  color: "#0D1030",
});

const SummaryCardBody = styled.div(
  {
    fontWeight: "bold",
    fontSize: "20px",
  },
  props => ({
    color: props.color ? props.color : "#3548D4",
  })
);

const FeaturedSummaryContainer = styled(Grid)({
  flexBasis: "60%",
});

const PieContainer = styled.div({
  display: "flex",
  justifyContent: "space-evenly",
});

const PieLegendContainer = styled.div({
  display: "flex",
  textAlign: "center",
  justifyContent: "space-evenly",
  flexDirection: "column",
});

interface FeaturedSummaryProps {
  name: string;
  data: {
    id: string;
    value: number;
    color: string;
  }[];
}

const FeaturedSummary = ({ summary }: { summary: FeaturedSummaryProps }) => {
  const values = summary?.data?.map(d => d.value);
  const total = values.reduce((t, value) => t + value);
  return (
    <FeaturedSummaryContainer item>
      <Paper>
        <SummaryCardTitle>{summary.name}</SummaryCardTitle>
        <PieContainer>
          <Pie
            data={summary?.data || []}
            height={215}
            width={215}
            innerRadius={0.75}
            colors={{ datum: "data.color" }}
            enableSliceLabels={false}
            enableRadialLabels={false}
            isInteractive={false}
          />
          <PieLegendContainer>
            <div>
              <SummaryCardTitle>Total</SummaryCardTitle>
              <SummaryCardBody>{total}</SummaryCardBody>
            </div>
            {summary?.data?.map(d => (
              <div key={d.id}>
                <SummaryCardTitle>{d.id}</SummaryCardTitle>
                <SummaryCardBody color={d.color}>{d.value}</SummaryCardBody>
              </div>
            ))}
          </PieLegendContainer>
        </PieContainer>
      </Paper>
    </FeaturedSummaryContainer>
  );
};

const SummariesContainer = styled(Grid)({
  textAlign: "center",
  flexBasis: "40%",
});

const InformationContainer = styled.div({
  padding: "16px 0",
});

interface DashboardProps {
  serverInfo: IClutch.envoytriage.v1.IServerInfo;
  featuredSummary: FeaturedSummaryProps;
  summaries?: {
    name: string;
    value: number;
  }[];
}

const Dashboard = ({ serverInfo, featuredSummary, summaries }: DashboardProps) => {
  const INFORMATION_KEYS = [
    "hot_restart_version",
    "uptime_all_epochs",
    "uptime_current_epoch",
    "version",
  ];

  const serverData = INFORMATION_KEYS.map(key => {
    return { name: key, value: serverInfo.value?.[key] };
  });

  return (
    <div>
      <Grid container direction="row" justify="space-evenly" wrap="nowrap" spacing={1}>
        <FeaturedSummary summary={featuredSummary} />
        <SummariesContainer item container direction="column" justify="space-evenly" spacing={1}>
          {summaries.map(summary => (
            <Grid item key={summary.name}>
              <Paper>
                <SummaryCardTitle>{summary.name}</SummaryCardTitle>
                <SummaryCardBody>{summary.value}</SummaryCardBody>
              </Paper>
            </Grid>
          ))}
        </SummariesContainer>
      </Grid>
      <InformationContainer>
        <MetadataTable data={serverData} />
      </InformationContainer>
    </div>
  );
};

export default Dashboard;
