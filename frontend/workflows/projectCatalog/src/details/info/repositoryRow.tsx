import React from "react";
import { Grid } from "@clutch-sh/core";
import { faCode, faCodeBranch } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

import { LinkText } from "../cards/base";

import type { ProjectRepository } from "./types";

const RepositoryRow = ({ name, url, icon = faCode, requests }: ProjectRepository) => (
  <>
    <Grid item>
      <FontAwesomeIcon icon={icon} size="lg" />
    </Grid>
    <Grid item>
      <LinkText text={name} link={url} />
    </Grid>
    {requests && (
      <>
        <Grid item>
          <FontAwesomeIcon icon={faCodeBranch} size="1x" />
        </Grid>
        <Grid item>
          <LinkText text={`${requests.number} ${requests.type}`} link={requests.url} />
        </Grid>
      </>
    )}
  </>
);

export default RepositoryRow;
