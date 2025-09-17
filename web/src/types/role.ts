
export type Role = {
    id: number;
    name: RoleName;
    description: string;
    createdAt: Date;
    updatedAt: Date;
}

export const RoleNameValue = {
    ADMIN: "Admin",
    DOCTOR: "Doctor",
    ANALYZER: "Analyst",
} as const


export type RoleName = typeof RoleNameValue[keyof typeof RoleNameValue]