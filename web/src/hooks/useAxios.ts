import axios, { type CreateAxiosDefaults } from "axios";
import { useMemo } from "react";
import { LOCAL_STORAGE_ACCESS_TOKEN } from "../types/constant";

/**
 * React hook to create and memoize an Axios instance.
 *
 * @param {object} [config] - Axios configuration object.
 * @returns {axios.AxiosInstance} - Memoized Axios instance.
 */
const useAxios = (config?: CreateAxiosDefaults) => {
  const instance = useMemo(() => {
    return createAxiosInstance(config);
  }, [config]); // Re-create instance only if config changes

  return instance;
};

/**
 * Creates an Axios instance with the given configuration
 * @param config - Optional Axios configuration object
 * @returns Configured Axios instance
 */
export const createAxiosInstance = (config?: CreateAxiosDefaults) => {
  const accessToken = localStorage.getItem(LOCAL_STORAGE_ACCESS_TOKEN);
  const conf = {
    baseURL: import.meta.env.VITE_BACKEND_BASE_URL,
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${accessToken}`,
    },
    ...config,
  };

  return axios.create(conf);
};

export default useAxios;