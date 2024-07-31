import React from "react";
import styled from "@emotion/styled";
import CloseIcon from "@mui/icons-material/Close";
import SearchIcon from "@mui/icons-material/Search";
import {
  alpha,
  ClickAwayListener,
  Grid,
  Icon,
  IconButton,
  InputAdornment as MuiInputAdornment,
  Popper as MuiPopper,
  TextField,
  TextFieldProps,
  Theme,
  Typography,
} from "@mui/material";
import type { AutocompleteRenderInputParams } from "@mui/material/Autocomplete";
import Autocomplete from "@mui/material/Autocomplete";
import type { FilterOptionsState } from "@mui/material/useAutocomplete";
import _ from "lodash";

import { useAppContext } from "../Contexts";
import { useNavigate } from "../navigation";

import type { SearchIndex } from "./utils";
import { filterHiddenRoutes, searchIndexes } from "./utils";

const hotKey = "/";

const InputField: React.FC<TextFieldProps> = styled(TextField)(({ theme }: { theme: Theme }) => ({
  // input field
  maxWidth: "551px",
  minWidth: "551px",
  "@media screen and (max-width: 880px)": {
    minWidth: "125px",
  },
  ".MuiInputBase-root": {
    height: "46px",
    border: `1px solid ${theme.palette.primary[600]}`,
    borderRadius: "4px",
    background: theme.palette.contrastColor,
    "&.Mui-focused fieldset": {
      border: "none",
    },
  },
  // input text color
  ".MuiAutocomplete-input": {
    color: theme.palette.secondary[900],
  },

  // close icon's container
  "div.MuiAutocomplete-endAdornment": {
    ".MuiAutocomplete-popupIndicatorOpen": {
      width: "32px",
      height: "32px",
      borderRadius: "30px",
      marginRight: "8px",
      "&:hover": {
        background: theme.palette.secondary[200],
      },
      "&:active": {
        background: theme.palette.secondary[300],
      },
    },
  },
}));

// search's result options container
const ResultGrid = styled(Grid)({
  height: "inherit",
  padding: "12px 16px 12px 16px",
});

// search's result options
const ResultLabel = styled(Typography)(({ theme }: { theme: Theme }) => ({
  color: theme.palette.secondary[900],
  fontSize: "14px",
}));

// main search icon on header
const SearchIconButton = styled(IconButton)(({ theme }: { theme: Theme }) => ({
  color: theme.palette.common.white,
  fontSize: "24px",
  padding: "12px",
  marginRight: "8px",
  "&:hover": {
    background: theme.palette.primary[600],
  },
  "&:active": {
    background: theme.palette.primary[700],
  },
}));

// search icon in input field
const StartInputAdornment = styled(MuiInputAdornment)(({ theme }: { theme: Theme }) => ({
  color: theme.palette.secondary[900],
  marginLeft: "8px",
}));

// closed icon svg
const StyledCloseIcon = styled(Icon)(({ theme }: { theme: Theme }) => ({
  color: theme.palette.secondary[900],
  fontSize: "24px",
}));

// popper containing the search result options
const Popper = styled(MuiPopper)(({ theme }: { theme: Theme }) => ({
  ".MuiPaper-root": {
    border: `1px solid ${theme.palette.secondary[100]}`,
    boxShadow: `0px 5px 15px ${alpha(theme.palette.primary[600], 0.2)}`,

    "> .MuiAutocomplete-listbox": {
      "> .MuiAutocomplete-option": {
        height: "48px",
        padding: "0px",

        "&.Mui-focused": {
          background: theme.palette.primary[200],
        },
      },
    },
  },
  ".MuiAutocomplete-noOptions": {
    fontSize: "14px",
    color: theme.palette.secondary[900],
  },
}));

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
  const { InputProps } = params;
  const searchRef = React.useRef();
  const handleKeyPress = (event: KeyboardEvent) => {
    if (searchRef.current) {
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
      autoFocus
      placeholder="Search..."
      fullWidth
      inputRef={searchRef}
      InputProps={{
        ...InputProps,
        disableUnderline: true,
        startAdornment: (
          <>
            <StartInputAdornment position="start">
              <SearchIcon />
            </StartInputAdornment>
            {InputProps.startAdornment}
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
  <ResultGrid container alignItems="center" onClick={handleSelection}>
    <Grid item xs>
      <ResultLabel>{option.label}</ResultLabel>
    </Grid>
  </ResultGrid>
);

const filterResults = (searchOptions: SearchIndex[], state: FilterOptionsState<SearchIndex>) =>
  _.filter(searchOptions, o => o.label.toLowerCase().includes(state.inputValue.toLowerCase()));

const SearchField: React.FC = () => {
  const { workflows } = useAppContext();
  const navigate = useNavigate();
  const options = searchIndexes(filterHiddenRoutes(workflows));
  const [inputValue, setInputValue] = React.useState("");
  const [showOptions, setShowOptions] = React.useState(false);
  const [open, setOpen] = React.useState(false);

  const renderResult = (props, option: SearchIndex) => {
    const handleSelection = () => {
      navigate(option.path);
    };

    return (
      <li {...props}>
        <Result option={option} handleSelection={handleSelection} />
      </li>
    );
  };

  const onInputChange = (__: React.ChangeEvent<{}>, value: string) => {
    if (value === "") {
      setShowOptions(false);
      setInputValue("");
    } else if (value !== hotKey) {
      setShowOptions(true);
      setInputValue(value);
    }

    // If full match will auto navigate user to the workflow
    const option = _.find(options, o => o.label.toLowerCase() === value.toLowerCase());
    if (option !== undefined) {
      setShowOptions(false);
      setInputValue("");
      navigate(option.path);
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

  // If workflow selected by pressing enter/return,
  // update the open state to collapse search bar to search icon
  const handleListKeyDown = event => {
    if (event.key === "Enter") {
      event.preventDefault();
      setOpen(false);
    }
  };

  // Will hide the search field if there are no visible workflows
  if (!workflows.length) {
    return null;
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
            getOptionLabel={x => (typeof x === "object" ? x.label : x)}
            PopperComponent={renderPopper}
            popupIcon={<CustomCloseIcon />}
            forcePopupIcon={!!showOptions}
            noOptionsText="No results found"
            onKeyDown={handleListKeyDown}
          />
        </ClickAwayListener>
      ) : (
        <SearchIconButton onClick={handleOpen} edge="end">
          <SearchIcon />
        </SearchIconButton>
      )}
    </Grid>
  );
};

export default SearchField;
