import type { URLSearchParamsInit } from "react-router-dom";
import {
  createSearchParams,
  useLocation,
  useNavigate as rrUseNavigate,
  useParams,
  useSearchParams as rrUseSearchParams,
} from "react-router-dom";
import type { PartialPath, State } from "history";

const UTM_PARAMS = ["utm_source", "utm_medium"];

/**
 * A custom partial Path object that may be missing some properties.
 *
 * This custom interface exists to provide an easier interface for
 * setting search params. See https://github.com/ReactTraining/react-router/issues/7743#issuecomment-785435977
 * for additional context.
 */
interface CustomPartialPath extends Pick<PartialPath, "pathname" | "hash"> {
  /**
   * Search params to include in navigation.
   */
  search?: URLSearchParamsInit;
}

type To = string | CustomPartialPath;

/**
 * React Router specific navigation options.
 */
interface ReactRouterNavigateOptions {
  replace?: boolean;
  state?: State;
}

/**
 * Custom navigation options to control Clutch specific routing functionality.
 */
export interface NavigateOptions extends ReactRouterNavigateOptions {
  // Include current URL in state as origin. Defaults to false.
  origin?: boolean;
  // Persist UTM search params after navigation. Defaults to true.
  utm?: boolean;
}

/**
 * Convert custom search param object to string.
 *
 * ```ts
 * const searchObject = new URLSearchParams({ query: "value" });
 * const searchString = convertSearchParam(searchObject);
 * console.log(searchString); // "?query=value"
 * ```
 */
const convertSearchParam = (params: URLSearchParams) => {
  const p = params.toString();
  if (p.indexOf("=") < 0) {
    return "";
  }
  return `?${p}`;
};

/**
 * Returns an imperative method for changing the location.
 *
 * This method wraps useNavigate from react-router but has custom
 * logic that:
 *   * Automatically preserves UTM parameters during navigation
 *   * Easily preserves the origin of navigation in state
 *   * Provides an easy interface to specify search parameters
 *
 * @see https://reactrouter.com/api/useNavigate for more information.
 */
const useNavigate = () => {
  const navigation = rrUseNavigate();
  const [currentSearchParams] = rrUseSearchParams();
  const customNavigate = (to: To, options: NavigateOptions = {}) => {
    const finalNavOptions = {} as ReactRouterNavigateOptions;
    if (options?.replace) {
      finalNavOptions.replace = options.replace;
    }
    if (options?.state) {
      finalNavOptions.state = options.state;
    }
    let newSearchParams = {} as URLSearchParamsInit;
    if (typeof to !== "string") {
      newSearchParams = (to?.search || {}) as URLSearchParamsInit;
    }
    const searchParams = createSearchParams(newSearchParams);
    if (options?.utm || options?.utm === undefined) {
      UTM_PARAMS.forEach(p => {
        const param = currentSearchParams.get(p);
        if (param && !searchParams.get(p)) {
          searchParams.set(p, param);
        }
      });
    }

    if (options?.origin) {
      const origin = `${window.location.pathname}${window.location.search}`;
      // n.b. if origin is specified in the options don't overwrite it.
      // @ts-ignore
      if (!options.state?.origin) {
        // @ts-ignore
        finalNavOptions.state = { ...finalNavOptions.state, origin };
      }
    }

    let navPath = {} as PartialPath;
    if (typeof to === "string") {
      navPath.pathname = to;
    } else {
      navPath = { pathname: to?.pathname, hash: to?.hash };
    }
    const [path, search] = navPath.pathname.split("?");
    navPath.pathname = path;
    if (search) {
      new URLSearchParams(search).forEach((v, k) => {
        searchParams.set(k, v);
      });
    }
    navPath.search = convertSearchParam(searchParams);
    navigation(navPath, finalNavOptions);
  };

  return customNavigate;
};

/**
 * Custom search param options to control Clutch specific routing functionality.
 */
interface SearchParamOptions extends ReactRouterNavigateOptions {
  utm?: boolean;
}

/**
 * A convienence wrapper for reading and writing search parameters via the
 * URLSearchParams interface. This custom hook wraps react-routers implementation
 * but changes the function to write search parameters to preserve UTM parameters
 * by default.
 */
const useSearchParams = (): readonly [
  URLSearchParams,
  (params: URLSearchParamsInit, options?: SearchParamOptions | undefined) => void
] => {
  const [searchParams, setSearchParams] = rrUseSearchParams();

  const customSetSearchParams = (
    params: URLSearchParamsInit,
    options?: SearchParamOptions | undefined
  ) => {
    const newSearchParams = params;
    const reactRouterOptions = {} as ReactRouterNavigateOptions;
    if (options?.replace) {
      reactRouterOptions.replace = options.replace;
    }
    if (options?.state) {
      reactRouterOptions.state = options.state;
    }

    if (options?.utm || options?.utm === undefined) {
      UTM_PARAMS.forEach(p => {
        const param = searchParams.get(p);
        if (param && !newSearchParams[p]) {
          newSearchParams[p] = param;
        }
      });
    }

    setSearchParams(newSearchParams, reactRouterOptions);
  };

  return [searchParams, customSetSearchParams];
};

export { convertSearchParam, useLocation, useParams, useSearchParams, useNavigate };
