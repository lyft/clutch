import React from "react";
import { FiSearch } from "react-icons/fi";
import { GrFormClose } from "react-icons/gr";
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
    border: "2px solid rgba(13, 16, 48, 0.1)",
    borderRadius: "4px",
    background: "white",
  },
  // input text color
  ".MuiAutocomplete-input": {
    color: "#0d1030",
  },

  // close icon's container
  "div.MuiAutocomplete-endAdornment":{
    ".MuiAutocomplete-popupIndicatorOpen": {
      width: "32px",
      height: "32px",
      borderRadius: "30px",
      marginRight: "8px",
      "&:hover": {
        background: "#e7e7ea",
      },
      "&:active": {
        background: "linear-gradient(0deg, rgba(255, 255, 255, 0.85), rgba(255, 255, 255, 0.85)), #0D1030;",
      },
    },
  },
});

// search's result options
const ResultLabel = styled(Typography)({
  color: "rgba(13, 16, 48, 0.6)",
  fontWeight: 500,
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
  fontSize: "24px",
});

// closed icon svg
const StyledCloseIcon = styled(Icon)({
    color: "#0c0b31",
    fontSize: "24px",
});

// popper with the search result options
const Popper = styled(MuiPopper)({
  ".MuiAutocomplete-paper": {
    border: "1px solid rgba(13, 16, 48, 0.12)",
    boxShadow: "0px 10px 24px rgba(35, 48, 143, 0.3)",
  },
  ".MuiAutocomplete-option": {
    height: "48px",
  },
  ".MuiAutocomplete-option[data-focus='true']": {
    background: "#ebedfb",
  },
});

const renderPopper = props => {
  return <Popper {...props} />;
};

const CustomCloseIcon: React.FC = () => {
  return (
      <StyledCloseIcon>
        <GrFormClose />
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
              <FiSearch />
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
  <Grid container alignItems="center" onClick={handleSelection}>
    <Grid item xs>
      <ResultLabel>{option.label}</ResultLabel>
    </Grid>
  </Grid>
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
            popupIcon={<CustomCloseIcon/>}
            forcePopupIcon={showOptions ? true: false}
            noOptionsText='No results found'
          />
        </ClickAwayListener>
      ) : (
        <SearchIconButton onClick={handleOpen}>
          <FiSearch />
        </SearchIconButton>
      )}
    </Grid>
  );
};

export default SearchField;
