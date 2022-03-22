import type { clutch as IClutch } from "@clutch-sh/api";

interface DeployEventList {
  events?: IClutch.timeseries.v1.IPoint[] | null;
}

// types for the Dash Deploys card
export type DeploysProjectMap = Record<string, DeployEventList>;
