import React, { useEffect, useState } from 'react'
import useAxios from '../hooks/useAxios';
// we intentionally avoid react-query here because this component is rendered
// outside of the application's QueryClientProvider. Use a direct axios call
// so the "What's New" modal can fetch server info without requiring react-query.

type Props = {
    children: React.ReactNode
}

const STORAGE_KEY = 'tamalabs_whats_new_seen_v1'

export default function WhatsNew({ children }: Props) {
    const axios = useAxios();
    const [serverInfo, setServerInfo] = useState<any | null>(null);

    useEffect(() => {
        let mounted = true;
        axios.get('/ping')
            .then((res) => {
                if (!mounted) return;
                setServerInfo(res.data);
            })
            .catch(() => {
                // ignore
            });

        return () => {
            mounted = false;
        };
    }, [axios]);

    const [seen, setSeen] = useState<boolean | null>(null)

    useEffect(() => {
        try {
            const v = localStorage.getItem(STORAGE_KEY)
            setSeen(v === 'true')
        } catch (e) {
            // If localStorage isn't available, assume seen to avoid blocking the app
            setSeen(true)
        }
    }, [])

    const markSeen = () => {
        try {
            localStorage.setItem(STORAGE_KEY, 'true')
        } catch (e) {
            // ignore
        }
        setSeen(true)
    }

    // still loading seen flag
    if (seen === null) return null

    if (seen) return <>{children}</>

    // render a simple full-screen modal
    return (
        <>
            <div style={{ position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.5)', zIndex: 999999999999, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                <div style={{ width: '90%', maxWidth: 800, background: 'white', borderRadius: 8, padding: 24, boxShadow: '0 8px 24px rgba(0,0,0,0.2)' }}>
                    <h2 style={{ marginTop: 0 }}>What's New?</h2>
                    <p>Welcome to the {serverInfo?.version ?? 'latest'} version of TAMALabs. Here are a few highlights:</p>
                    <ul>
                        <li>Added an interactive Dashboard for better data visualization and quick insights.</li>
                        <li>Introduced License Management feature to simplify user and license handling.</li>
                        <li>Improved print layout for reports and documents for clearer and more professional output.</li>
                        <li>Fixed several bugs to enhance overall system stability and performance.</li>
                        <li>Enhanced user interface for a smoother and more intuitive experience.</li>
                        <li>Optimized performance for faster load times and responsiveness.</li>
                        <li>Improved authentication system for more secure and reliable user access.</li>
                        <li>Automatic database backup to ensure data safety and easy recovery.</li>
                    </ul>
                    <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8, marginTop: 16 }}>
                        <button onClick={markSeen} style={{ padding: '8px 12px', borderRadius: 4, border: '1px solid #ccc', background: '#fff', cursor: 'pointer' }}>Got it</button>
                    </div>
                </div>
            </div>
            {/* children still render behind the modal so app state initializes; modal blocks interaction */}
            <>{children}</>
        </>
    )
}


