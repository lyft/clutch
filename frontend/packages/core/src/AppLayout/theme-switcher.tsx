import React from "react";
import DarkModeIcon from "@mui/icons-material/DarkMode";
import LightModeIcon from "@mui/icons-material/LightMode";
import { Grid } from "@mui/material";
import get from "lodash/get";

import { useUserPreferences } from "../Contexts";
import { Select } from "../Input";
import { THEME_VARIANTS } from "../Theme/colors";

const ThemeSwitcher = () => {
  const { preferences, dispatch } = useUserPreferences();

  const options = [
    {
      label: THEME_VARIANTS.light,
      startAdornment: <LightModeIcon />,
    },
    {
      label: THEME_VARIANTS.dark,
      startAdornment: <DarkModeIcon />,
    },
  ];

  const handleOnChange = (value: string) => {
    dispatch({
      type: "SetPref",
      payload: { key: "theme", value },
    });
  };

  return (
    <Grid container alignItems="center" justifyContent="space-between">
      <Grid item>Theme</Grid>
      <Grid item>
        <Select
          name="theme-switcher"
          label=""
          options={options}
          onChange={handleOnChange}
          defaultOption={get(preferences, "theme")}
        />
      </Grid>
    </Grid>
  );
};

export default ThemeSwitcher;
