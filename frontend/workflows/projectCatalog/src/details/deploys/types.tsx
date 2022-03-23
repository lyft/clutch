export type Environments = "SETUP" | "STAGING" | "CANARY" | "PRODUCTION";

interface CommitAuthor {
  username?: string | null;
  email?: string | null;
}

export interface Commit {
  /** Commit ref */
  ref?: string | null;

  /** Commit message */
  message?: string | null;

  /** Commit author */
  author?: CommitAuthor | null;
}

export interface CommitInfo {
  /** The name of repository including owner / org */
  repositoryName: string;
  /** The commits that will be spanned for the comparison */
  commits: Commit[];
  /** The base ref that will be used for comparisons. This should be one commit before the first one in commits. */
  baseRef?: string;
  url?: string;
  environment?: Environments;
}

export interface ProjectDeploys {
  title?: string;
  lastDeploy?: number;
  deploys: CommitInfo[];
  seeMore: {
    text: string;
    url: string;
  };
}
