import React from "react";
import { useNavigate } from "react-router-dom";
import styled from "@emotion/styled";
import {
  ClickAwayListener,
  Grid,
  Icon,
  IconButton,
  InputAdornment as MuiInputAdornment,
  Popper as MuiPopper,
  TextField,
  Typography,
} from "@material-ui/core";
import CloseIcon from "@material-ui/icons/Close";
import SearchIcon from "@material-ui/icons/Search";
import type { AutocompleteRenderInputParams } from "@material-ui/lab/Autocomplete";
import Autocomplete from "@material-ui/lab/Autocomplete";
import type { FilterOptionsState } from "@material-ui/lab/useAutocomplete";
import _ from "lodash";

import { useAppContext } from "../Contexts";

import type { SearchIndex } from "./utils";
import { searchIndexes } from "./utils";

const hotKey = "/";

const InputField = styled(TextField)({
  // input field
  maxWidth: "551px",
  minWidth: "551px",
  "@media screen and (max-width: 880px)": {
    minWidth: "125px",
  },
  ".MuiInputBase-root": {
    height: "46px",
    border: "1px solid #3548d4",
    borderRadius: "4px",
    background: "#ffffff",
  },
  // input text color
  ".MuiAutocomplete-input": {
    color: "#0d1030",
  },

  // close icon's container
  "div.MuiAutocomplete-endAdornment": {
    ".MuiAutocomplete-popupIndicatorOpen": {
      width: "32px",
      height: "32px",
      borderRadius: "30px",
      marginRight: "8px",
      "&:hover": {
        background: "#e7e7ea",
      },
      "&:active": {
        background: "#DBDBE0",
      },
    },
  },
});

const ResultGrid = styled(Grid)({
  height: "inherit",
})

// search's result options
const ResultLabel = styled(Typography)({
  color: "#0d1030",
  fontSize: "14px",
});

// main search icon on header
const SearchIconButton = styled(IconButton)({
  color: "#ffffff",
  fontSize: "24px",
  padding: "12px",
  "&:hover": {
    background: "#2d3db4",
  },
  "&:active": {
    background: "#2938a5",
  },
});

// search icon in input field
const StartInputAdornment = styled(MuiInputAdornment)({
  color: "#0c0b31",
  marginLeft: "8px",
});

// closed icon svg
const StyledCloseIcon = styled(Icon)({
  color: "#0c0b31",
  fontSize: "24px",
});

// popper containing the search result options
const Popper = styled(MuiPopper)({
  ".MuiAutocomplete-paper": {
    border: "1px solid #e7e7ea",
    boxShadow: "0px 5px 15px rgba(53, 72, 212, 0.2)",
  },
  ".MuiAutocomplete-option": {
    height: "48px",
  },
  ".MuiAutocomplete-option[data-focus='true']": {
    background: "#ebedfb",
  },
  ".MuiAutocomplete-noOptions": {
    fontSize: "14px",
    color: "#0d1030",
  },
});

const renderPopper = props => {
  return <Popper {...props} />;
};

const CustomCloseIcon: React.FC = () => {
  return (
    <StyledCloseIcon>
      <CloseIcon fontSize="small" />
    </StyledCloseIcon>
  );
};

const Input = (params: AutocompleteRenderInputParams): React.ReactNode => {
  const searchRef = React.useRef();
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

  return (
    <InputField
      {...params}
      placeholder="Search..."
      fullWidth
      inputRef={searchRef}
      InputProps={{
        ...params.InputProps,
        disableUnderline: true,
        startAdornment: (
          <>
            <StartInputAdornment position="start">
              <SearchIcon />
            </StartInputAdornment>
            {params.InputProps.startAdornment}
          </>
        ),
      }}
    />
  );
};

interface ResultProps {
  option: SearchIndex;
  handleSelection: () => void;
}

const Result: React.FC<ResultProps> = ({ option, handleSelection }) => (
  <ResultGrid container alignItems="center" onClick={handleSelection} style={{height: "inherit"}}>
    <Grid item xs>
      <ResultLabel>{option.label}</ResultLabel>
    </Grid>
  </ResultGrid>
);

const filterResults = (searchOptions: SearchIndex[], state: FilterOptionsState<SearchIndex>) => {
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
  const [open, setOpen] = React.useState(false);

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

  const handleOpen = () => {
    setOpen(!open);
  };

  const handleClose = () => {
    setOpen(false);
  };

  function handleListKeyDown(event) {
    if (event.key === "Enter") {
      event.preventDefault();
      setOpen(false);
    }
  }

  return (
    <Grid container alignItems="center">
      {open ? (
        <ClickAwayListener onClickAway={handleClose}>
          <Autocomplete
            autoComplete
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
            PopperComponent={renderPopper}
            popupIcon={<CustomCloseIcon />}
            forcePopupIcon={!!showOptions}
            noOptionsText="No results found"
            onKeyDown={handleListKeyDown}
          />
        </ClickAwayListener>
      ) : (
        <SearchIconButton onClick={handleOpen}>
          <SearchIcon />
        </SearchIconButton>
      )}
    </Grid>
  );
};

export default SearchField;
