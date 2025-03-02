import axios, { type CreateAxiosDefaults } from "axios";
import { useMemo } from "react";

/**
 * React hook to create and memoize an Axios instance.
 *
 * @param {object} [config] - Axios configuration object.
 * @returns {axios.AxiosInstance} - Memoized Axios instance.
 */
const useAxios = (config?: CreateAxiosDefaults) => {
  const instance = useMemo(() => {
    const conf = {
      baseURL: import.meta.env.VITE_BACKEND_BASE_URL,
      ...config,
    };

    return axios.create(conf);
  }, [config]); // Re-create instance only if config changes

  return instance;
};

export default useAxios;
