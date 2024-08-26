import type { AlertProps as MuiAlertProps } from "@mui/lab";

export interface AlertProps
  extends Pick<
    MuiAlertProps,
    "severity" | "action" | "onClose" | "elevation" | "variant" | "icon" | "className"
  > {
  title?: React.ReactNode;
}

export interface Banner extends Pick<AlertProps, "title" | "severity"> {
  message: string;
  dismissed: boolean;
  linkText?: string;
  link?: string;
}

export interface PerWorkflowBanner {
  [workflowName: string]: Banner;
}

export interface WorkflowsBanner extends Banner {
  workflows: string[];
}

export interface AppBanners {
  /** Will display a notification banner at the top of the application */
  header?: Banner;
  /** Allows for setting a notification banner on a per workflow basis */
  perWorkflow?: PerWorkflowBanner;
  /** Allows for setting a notification banner across multiple workflows */
  multiWorkflow?: WorkflowsBanner;
}
