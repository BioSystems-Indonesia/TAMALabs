import { LOCAL_STORAGE_ADMIN } from "../types/constant";
import { User } from "../types/user";

export const useCurrentUser = () => {
  const admin = localStorage.getItem(LOCAL_STORAGE_ADMIN);
  if (!admin) {
    return null;
  }

  return JSON.parse(admin) as User;
};

export const useCurrentUserRole = () => {
  // Get role from localStorage admin data instead of JWT token
  const admin = localStorage.getItem(LOCAL_STORAGE_ADMIN);
  if (!admin) {
    throw new Error("No admin data found");
  }

  try {
    return JSON.parse(admin).roles[0].name;
  } catch (error) {
    throw new Error("Invalid admin data");
  }
};
