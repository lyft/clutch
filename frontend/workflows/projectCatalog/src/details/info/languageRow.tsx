import React from "react";
import { Grid, Typography } from "@clutch-sh/core";

import LanguageIcon from "../../helpers/language-icon";

const LanguageRow = ({ languages }: { languages: string[] }) => (
  <>
    <Grid item>
      <Typography variant="body2">Language{languages.length > 1 && "s"}</Typography>
    </Grid>
    {languages.map(language => (
      <Grid item>
        <LanguageIcon language={language} />
      </Grid>
    ))}
  </>
);

export default LanguageRow;
