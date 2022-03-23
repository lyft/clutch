import type { clutch as IClutch } from "@clutch-sh/api";
import type { ClutchError } from "@clutch-sh/core";

import type { AlertInfo } from "./alertsOld/types";
import type { DeployInfo } from "./deploysOld/types";
import type { ProjectInfo } from "./info/types";

type DeploysActionKind = "DEPLOYS_START" | "DEPLOYS_STOP" | "DEPLOYS_HYDRATE" | "DEPLOYS_ERROR";
type AlertsActionKind = "ALERTS_START" | "ALERTS_STOP" | "ALERTS_HYDRATE" | "ALERTS_ERROR";
type InfoActionKind = "INFO_START" | "INFO_STOP" | "INFO_HYDRATE" | "INFO_ERROR";

interface DeploysPayload {
  result: DeployInfo;
}
interface AlertsPayload {
  result: AlertInfo;
}
interface InfoPayload {
  result: ProjectInfo;
}

interface DeploysAction {
  type: DeploysActionKind;
  payload?: DeploysPayload;
}

interface AlertsAction {
  type: AlertsActionKind;
  payload?: AlertsPayload;
}

interface InfoAction {
  type: InfoActionKind;
  payload?: InfoPayload;
}

export type Action = InfoAction | AlertsAction | DeploysAction;

interface DefaultState {
  error?: ClutchError | undefined;
  loading?: boolean;
}

interface AlertsState extends DefaultState {
  data?: AlertInfo;
}

interface DeploysState extends DefaultState {
  data?: DeployInfo;
}

interface InfoState extends DefaultState {
  data?: ProjectInfo;
}

export interface DetailsState {
  alerts: AlertsState;
  deploys: DeploysState;
  info: InfoState;
}

interface DeployEventList {
  events?: IClutch.timeseries.v1.IPoint[] | null;
}

// types for the Dash Deploys card
export type DeploysProjectMap = Record<string, DeployEventList>;
