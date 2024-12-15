export const Action = {
    CREATE: 'CREATE',
    EDIT: 'EDIT',
    SHOW: 'SHOW'
} as const;

export type ActionKeys = keyof typeof Action;