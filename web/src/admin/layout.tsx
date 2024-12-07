import type { ReactNode } from 'react';
import { CheckForApplicationUpdate, Layout } from 'react-admin';

export const DefaultLayout = ({ children }: { children: ReactNode }) => (
    <Layout>
        {children}
        <CheckForApplicationUpdate />
    </Layout>
);