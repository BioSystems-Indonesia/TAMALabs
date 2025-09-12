import { jwtDecode } from "jwt-decode";
import { LOCAL_STORAGE_ACCESS_TOKEN, LOCAL_STORAGE_ADMIN } from "../types/constant";
import { User } from "../types/user";

export const useCurrentUser = () => {
    const admin = localStorage.getItem(
        LOCAL_STORAGE_ADMIN,
    );
    if (!admin) {
        return null;
    }

    return JSON.parse(admin) as User;
}

export const useCurrentUserRole = () => {
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
};
