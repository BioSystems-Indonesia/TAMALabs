import type { AuthProvider, UserIdentity } from "react-admin";
import type { User } from "../types/user";
import { LOCAL_STORAGE_ADMIN } from "../types/constant";
import { createAxiosInstance } from "./useAxios";
import { AxiosError } from "axios";

function redirectToLogin(): never {
  try {
    localStorage.removeItem(LOCAL_STORAGE_ADMIN);
  } catch (e) {
    // ignore
  }
  window.location.href = "/#/login";
  throw new Error("Not authenticated");
}

const authProvider: AuthProvider = {
  async login({ username, password }) {
    try {
      const axios = createAxiosInstance();
      const response = await axios.post("/login", {
        username: username,
        password: password,
      });
      // Simpan data admin di localStorage untuk getIdentity()
      localStorage.setItem(
        LOCAL_STORAGE_ADMIN,
        JSON.stringify(response.data.admin)
      );
      return Promise.resolve();
    } catch (error) {
      if (error instanceof AxiosError) {
        if (error.response) {
          if (error.response.status >= 400 && error.response.status < 500) {
            throw new Error(
              error.response.data.error ?? "Invalid username or password"
            );
          } else if (error.response.status >= 500) {
            throw new Error("Server error");
          }
        }
        throw new Error("Network error");
      }
      throw new Error("An unexpected error occurred");
    }
  },
  async checkError(error) {
    const status = error.status;
    if (status === 401) {
      return redirectToLogin();
    }
  },
  async checkAuth() {
    try {
      const axios = createAxiosInstance();
      await axios.get("/check-auth");
    } catch (error) {
      if (error instanceof AxiosError && error.response?.status === 401) {
        return redirectToLogin();
      }
      throw error;
    }
  },
  async logout() {
    try {
      const axios = createAxiosInstance();
      await axios.post("/logout", {});
    } catch (error) {
      // Ignore logout errors, just redirect to login
      console.warn("Logout error:", error);
    }
    // Hapus data admin dari localStorage
    localStorage.removeItem(LOCAL_STORAGE_ADMIN);
    // Redirect to login page
    window.location.href = "/#/login";
  },
  async getIdentity(): Promise<UserIdentity> {
    const adminStorage = localStorage.getItem(LOCAL_STORAGE_ADMIN);
    if (!adminStorage) {
      return redirectToLogin();
    }

    try {
      const admin: User = JSON.parse(adminStorage);
      return {
        id: admin.id,
        fullName: admin.fullname,
      };
    } catch (error) {
      // Jika data di localStorage rusak, redirect ke login
      return redirectToLogin();
    }
  },
  async getPermissions(): Promise<string> {
    try {
      const axios = createAxiosInstance();
      const response = await axios.get("/permissions");
      return response.data.role || response.data.permission || "user";
    } catch (error) {
      if (error instanceof AxiosError && error.response?.status === 401) {
        return redirectToLogin();
      }

      // Fallback: return default role instead of throwing error
      console.warn(
        "Failed to get permissions from server, using default role:",
        error
      );
      return "user";
    }
  },
};

export const useAuthProvider = () => {
  return authProvider;
};
