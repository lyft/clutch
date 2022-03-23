import type { CommitInfo } from "./commitInformation";

export type Environments = "SETUP" | "STAGING" | "CANARY" | "PRODUCTION";

export type Statuses =
  | "UNKNOWN"
  | "WAITING"
  | "RUNNING"
  | "SUCCESS"
  | "FAILURE"
  | "ABORTED"
  | "SKIPPED"
  | "QUEUED";

export interface DeployJobInformation {
  name?: string;
  commitMetadata: CommitInfo;
  status: Statuses;
  timestamp: number;
  environment: Environments;
}

export interface DeployInfo {
  jobs: DeployJobInformation[];
  inProgress?: number;
  failures?: number;
  lastDeploy?: number | Long | null;
}
