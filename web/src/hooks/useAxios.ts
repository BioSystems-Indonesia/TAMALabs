import axios, { type CreateAxiosDefaults } from "axios";
import { useMemo } from "react";
import { useNotify } from "react-admin";
import { ErrorPayload } from "../types/errors";

/**
 * React hook to create and memoize an Axios instance.
 *
 * @param {object} [config] - Axios configuration object.
 * @returns {axios.AxiosInstance} - Memoized Axios instance.
 */
const useAxios = (config?: CreateAxiosDefaults) => {
  const notify = useNotify();
  const instance = useMemo(() => {
    return createAxiosInstance(config);
  }, [config]); // Re-create instance only if config changes

  instance.interceptors.response.use(
    (response) => response,
    (error) => {
      if (axios.isAxiosError(error)) {
        const errorResponse = error.response?.data as ErrorPayload;

        // If no error response, show generic error
        if (!errorResponse) {
          notify("Something went wrong", {
            type: "error",
          });
          return Promise.reject(error);
        }

        // Handle specific status codes with default notifications
        switch (errorResponse.status_code) {
          case 401:
            notify("Session expired", {
              type: "error",
            });
            break;
          case 403:
            break;
          default:
            if (!error.config?.url?.includes('/nuha-simrs/')) {
              notify(
                `Error ${errorResponse.status_code}: ${errorResponse.error}`,
                {
                  type: "error",
                }
              );
            }
            break;
        }

        // Always re-throw error so mutation onError can handle it
        return Promise.reject(error);
      } else {
        return Promise.reject(error);
      }
    }
  );

  return instance;
};

/**
 * Creates an Axios instance with the given configuration
 * @param config - Optional Axios configuration object
 * @returns Configured Axios instance
 */
export const createAxiosInstance = (config?: CreateAxiosDefaults) => {
  const conf = {
    baseURL: import.meta.env.VITE_BACKEND_BASE_URL,
    headers: {
      "Content-Type": "application/json",
    },
    withCredentials: true,
    ...config,
  };

  const instance = axios.create(conf);
  return instance;
};

export default useAxios;
