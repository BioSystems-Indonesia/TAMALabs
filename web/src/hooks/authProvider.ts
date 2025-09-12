import type {
  AuthProvider,
  UserIdentity
} from "react-admin";
import type { User } from "../types/user";
import {
  LOCAL_STORAGE_ACCESS_TOKEN,
  LOCAL_STORAGE_ADMIN,
} from "../types/constant";
import { createAxiosInstance } from "./useAxios";
import { AxiosError } from "axios";
import { jwtDecode } from "jwt-decode";

const authProvider: AuthProvider = {
  async login({ username, password }) {
    try {
      const axios = createAxiosInstance();
      const response = await axios.post("/login", {
        username: username,
        password: password,
      });
      localStorage.setItem(
        LOCAL_STORAGE_ACCESS_TOKEN,
        response.data.access_token
      );
      localStorage.setItem(
        LOCAL_STORAGE_ADMIN,
        JSON.stringify(response.data.admin)
      );

    } catch (error) {
      if (error instanceof AxiosError) {
        if (error.response) {
          if (error.response.status >= 400 && error.response.status < 500) {
            throw new Error(error.response.data.error ?? "Invalid username or password");
          } else if (error.response.status >= 500) {
            throw new Error("Server error");
          }
        }
        throw new Error('Network error');
      }
      throw new Error('An unexpected error occurred');
    }


  },
  async checkError(error) {
    const status = error.status;
    if (status === 401) {
      localStorage.removeItem(LOCAL_STORAGE_ACCESS_TOKEN);
      localStorage.removeItem(LOCAL_STORAGE_ADMIN);
      throw new Error("Session expired");
    }
  },
  async checkAuth() {
    try {
      const axios = createAxiosInstance();
      await axios.get("/check-auth");
    } catch (error) {
      if (error instanceof AxiosError) {
        if (error.response) {
          if (error.response.status === 401) {
            throw new Error("Session expired");
          }
        }
      }
    }
  },
  async logout() {
    localStorage.removeItem(LOCAL_STORAGE_ACCESS_TOKEN);
    localStorage.removeItem(LOCAL_STORAGE_ADMIN);
  },
  async getIdentity(): Promise<UserIdentity> {
    const adminStorage = localStorage.getItem(LOCAL_STORAGE_ADMIN);
    if (!adminStorage) {
      throw new Error("Not authenticated");
    }

    const admin: User = JSON.parse(adminStorage);
    return {
      id: admin.id,
      fullName: admin.fullname,
    };
  },
  async getPermissions(): Promise<string> {
    const token = localStorage.getItem(LOCAL_STORAGE_ACCESS_TOKEN);
    if (!token) {
      throw new Error("No token");
    }

    try {
      const decoded = jwtDecode(token) as { role: string };
      return decoded.role;
    } catch (error) {
      throw new Error("Invalid token");
    }
  },
};


export const useAuthProvider = () => {
  return authProvider;
}

