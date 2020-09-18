import React from "react";
import { useNavigate } from "react-router-dom";
import { Box, Grid, TextField, Typography, useTheme } from "@material-ui/core";
import { fade } from "@material-ui/core/styles";
import SearchIcon from "@material-ui/icons/Search";
import type { RenderInputParams } from "@material-ui/lab/Autocomplete";
import Autocomplete from "@material-ui/lab/Autocomplete";
import type { FilterOptionsState } from "@material-ui/lab/useAutocomplete";
import _ from "lodash";
import styled from "styled-components";

import { useAppContext } from "../Contexts";

import type { SearchIndex } from "./utils";
import { searchIndexes } from "./utils";

const hotKey = "/";

const Container = styled(Box)`
  ${({ theme }) => `
  border-radius: ${theme.shape.borderRadius}px;
  margin-left: 0px;
  margin-right: ${theme.spacing(2)}px;
  background-color: ${fade(theme.palette.primary.main, 0.15)};
  &:hover {
    background-color: ${fade(theme.palette.primary.main, 0.25)};
  }
  `}
`;

const InputField = styled(TextField)`
  max-width: 300px;
  min-width: 300px;
  @media screen and (max-width: 800px) {
    min-width: 125px;
  }
`;

const InputIcon = styled(SearchIcon)`
  margin-top: -3px;
`;

const ResultLabel = styled(Typography)`
  ${({ theme }) => `
  color: ${theme.palette.primary.main};
  font-weight: 400;
  `}
`;

const Input = (params: RenderInputParams): React.ReactNode => {
  const searchRef = React.useRef();
  const theme = useTheme();
  const handleKeyPress = (event: KeyboardEvent) => {
    if (searchRef.current !== undefined) {
      if (event.key === hotKey && (event.target as Node).nodeName !== "INPUT") {
        // @ts-ignore
        searchRef.current.focus();
      } else if (event.key === "Escape") {
        // @ts-ignore
        searchRef.current.blur();
      }
    }
  };

  React.useLayoutEffect(() => {
    window.addEventListener("keydown", handleKeyPress);
  }, []);

  const inputProps = { ...params.InputProps, style: { color: theme.palette.primary.main } };

  return (
    <InputField
      {...params}
      label={<InputIcon />}
      placeholder="Search..."
      variant="outlined"
      fullWidth
      inputRef={searchRef}
      InputProps={inputProps}
    />
  );
};

interface ResultProps {
  option: SearchIndex;
  handleSelection: () => void;
}

const Result: React.FC<ResultProps> = ({ option, handleSelection }) => (
  <Grid container alignItems="center" onClick={handleSelection}>
    <Grid item xs>
      <ResultLabel>{option.label}</ResultLabel>
      <Typography variant="body2" color="textSecondary">
        {option.category}
      </Typography>
    </Grid>
  </Grid>
);

const filterResults = (searchOptions: SearchIndex[], state: FilterOptionsState) => {
  return _.filter(searchOptions, o => {
    return o.label.toLowerCase().includes(state.inputValue.toLowerCase());
  });
};

const SearchField: React.FC = () => {
  const { workflows } = useAppContext();
  const navigate = useNavigate();
  const options = searchIndexes(workflows);
  const [inputValue, setInputValue] = React.useState("");
  const [showOptions, setShowOptions] = React.useState(false);
  const theme = useTheme();

  const renderResult = (option: SearchIndex) => {
    const handleSelection = () => {
      navigate(option.path);
    };

    return <Result option={option} handleSelection={handleSelection} />;
  };

  const onInputChange = (__: React.ChangeEvent<{}>, value: string) => {
    if (value === "") {
      setShowOptions(false);
      setInputValue("");
    } else if (value !== hotKey) {
      setShowOptions(true);
      setInputValue(value);
    }
    const option = _.find(options, o => {
      return o.label === value;
    });
    if (option !== undefined) {
      setShowOptions(false);
      setInputValue("");
    }
  };

  const onOptionsOpen = () => {
    if (inputValue !== "") {
      setShowOptions(true);
    }
  };

  const onOptionsClose = (event: any) => {
    setShowOptions(false);
    const option = _.find(options, o => {
      return o.label === event.target.value;
    });
    if (option !== undefined) {
      navigate(option.path);
    }
    setInputValue("");
  };

  return (
    <Container>
      <Autocomplete
        autoComplete
        freeSolo
        selectOnFocus
        size="small"
        inputValue={inputValue}
        renderInput={Input}
        renderOption={renderResult}
        onInputChange={onInputChange}
        open={showOptions}
        onOpen={onOptionsOpen}
        onClose={onOptionsClose}
        options={options}
        filterOptions={filterResults}
        getOptionLabel={x => x.label}
        ListboxProps={{ style: { backgroundColor: fade(theme.palette.secondary.main, 0.75) } }}
      />
    </Container>
  );
};

export default SearchField;
