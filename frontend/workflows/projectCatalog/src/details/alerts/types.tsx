import type { IconProp } from "@fortawesome/fontawesome-svg-core";

export interface AlertMassageOptions {
  offset?: number;
  statuses?: string[];
  title?: string;
  text?: string;
  icon?: IconProp;
  url?: string;
}

export interface User {
  name: string;
  url?: string;
}

export interface OnCall {
  text: string;
  icon?: IconProp;
  users?: User[];
  url?: string;
}

interface Summary {
  count: number;
  url?: string;
}

export interface AlertSummary {
  open?: Summary;
  triggered?: Summary;
  acknowledged?: Summary;
}

export interface ProjectAlerts {
  title?: string;
  lastAlert?: number;
  summary?: AlertSummary;
  onCall?: OnCall;
  create?: {
    text: string;
    url: string;
  };
}
