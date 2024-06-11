import React from "react";
import DarkModeIcon from "@mui/icons-material/DarkMode";
import LightModeIcon from "@mui/icons-material/LightMode";
import { Grid } from "@mui/material";
import get from "lodash/get";

import { useUserPreferences } from "../Contexts";
import { Select } from "../Input";

const ThemeSwitcher = () => {
  const { preferences, dispatch } = useUserPreferences();

  const options = [
    {
      label: "Light",
      startAdornment: <LightModeIcon />,
    },
    {
      label: "Dark",
      startAdornment: <DarkModeIcon />,
    },
  ];

  const handleOnChange = (value: string) => {
    dispatch({
      type: "SetPref",
      payload: { key: "theme", value: value.toUpperCase() },
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
