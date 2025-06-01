import { LOCAL_STORAGE_ADMIN } from "../types/constant";
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