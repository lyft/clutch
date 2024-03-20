import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";

import AuditEvent from "./audit-event";
import AuditLog from "@logs";

export interface AuditLogProps extends BaseWorkflowProps {
  detailsPathPrefix?: string;
  downloadPrefix?: string;
}

export interface WorkflowProps extends BaseWorkflowProps {}

const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "audit",
    group: "Audit",
    displayName: "Audit Trail",
    routes: {
      landing: {
        path: "/",
        displayName: "Logs",
        description: "View audit log",
        component: AuditLog,
      },
      event: {
        path: "/event/:id",
        displayName: "Event Details",
        description: "View audit event",
        component: AuditEvent,
        hideNav: true,
      },
    },
  };
};

export default register;
