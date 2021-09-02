import * as React from "react";
import { Table, TableRow, Typography } from "@clutch-sh/core";
import styled from "@emotion/styled";
import { Box, Grid as MuiGrid } from "@material-ui/core";
import _ from "lodash";

import Card from "./card";
import { DashDispatchContext, DashStateContext } from "./dash-hooks";
import ProjectSelector from "./project-selector";
import type { DashAction, DashState } from "./types";

const initialState = {
  selected: [],
  projectData: {},
};

const BigGrid = styled(MuiGrid)({
  margin: "7px",
});

const dashReducer = (state: DashState, action: DashAction): DashState => {
  switch (action.type) {
    case "UPDATE_SELECTED": {
      if (!_.isEqual(state.selected, action.payload.selected)) {
        return action.payload;
      }
      return state;
    }
    default:
      throw new Error("not implemented (should be unreachable)");
  }
};

const Dash = ({ children }) => {
  const [state, dispatch] = React.useReducer(dashReducer, initialState);

  return (
    <Box display="flex" flex={1} minHeight="100%">
      <DashDispatchContext.Provider value={dispatch}>
        <DashStateContext.Provider value={state}>
          <ProjectSelector />
          <Box display="flex" flex={1}>
            <BigGrid
              spacing={3}
              container
              direction="row"
              justify="flex-start"
              alignItems="flex-start"
              alignContent="flex-start"
            >
              <Card
                avatar="🚀"
                title="Deploys"
                summary={[
                  {
                    title: <Typography variant="subtitle2">5m ago</Typography>,
                    subheader: "Last deploy",
                  },
                  {
                    title: (
                      <Typography variant="subtitle2" color="#3548D4">
                        0
                      </Typography>
                    ),
                  },
                  {
                    title: (
                      <Typography variant="subtitle2" color="#DB3615">
                        0
                      </Typography>
                    ),
                    subheader: "Failed",
                  },
                ]}
              >
                <Table columns={["", "", ""]}>
                  <TableRow>
                    <div>clutch</div>
                    <div>No commits</div>
                    <div>🥚</div>
                  </TableRow>
                  <TableRow>
                    <div>clutch</div>
                    <div>No commits</div>
                    <div>🥚</div>
                  </TableRow>
                  <TableRow>
                    <div>clutch</div>
                    <div>No commits</div>
                    <div>🥚</div>
                  </TableRow>
                </Table>
              </Card>
              <Card
                avatar="🚀"
                title="Deploys"
                summary={[
                  {
                    title: <Typography variant="subtitle2">5m ago</Typography>,
                    subheader: "Last deploy",
                  },
                  {
                    title: (
                      <Typography variant="subtitle2" color="#3548D4">
                        3
                      </Typography>
                    ),
                    subheader: "In progress",
                  },
                  {
                    title: (
                      <Typography variant="subtitle2" color="#DB3615">
                        0
                      </Typography>
                    ),
                    subheader: "Failed",
                  },
                  {
                    title: <Typography variant="subtitle2">0</Typography>,
                    subheader: "Rollbacks",
                  },
                ]}
              >
                <Table columns={["", "", "", "", "", ""]}>
                  <TableRow>
                    <div>clutch</div>
                    <div>Commit</div>
                    <div>🥚</div>
                    <div>RUNNING</div>
                    <div>0 days ago</div>
                    <div>Initiated by ANON</div>
                  </TableRow>
                  <TableRow>
                    <div>clutch</div>
                    <div>Commit</div>
                    <div>🥚</div>
                    <div>RUNNING</div>
                    <div>0 days ago</div>
                    <div>Initiated by ANON</div>
                  </TableRow>
                  <TableRow>
                    <div>clutch</div>
                    <div>A longer commit message than other rows</div>
                    <div>🥚</div>
                    <div>RUNNING</div>
                    <div>0 days ago</div>
                    <div>Initiated by ANON</div>
                  </TableRow>
                </Table>
              </Card>
            </BigGrid>
          </Box>
        </DashStateContext.Provider>
      </DashDispatchContext.Provider>
    </Box>
  );
};

export default Dash;
