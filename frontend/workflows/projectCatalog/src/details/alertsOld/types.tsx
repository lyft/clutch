export type Urgency = "HIGH" | "LOW";
export type Status = "ack" | "open" | "closed";

export interface Assignment {
  assignee: string;
  at: string;
}

export interface Incident {
  id: string;
  urgency: Urgency;
  url?: string;
  created: string;
  description: string;
  assignments?: Assignment[];
  status: Status;
}

export interface Alert {
  incidents?: Incident[];
}

export interface Alerts {
  [key: string]: Alert;
}

export interface AlertInfo {
  lastAlert?: number | Long | null;
  acknowledged?: number;
  open?: number;
  projectAlerts?: Alerts;
}
