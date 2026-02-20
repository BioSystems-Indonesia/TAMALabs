import { createContext, useContext, useEffect, useRef, useState, type ReactNode } from "react";

type ServiceStatus = "connected" | "disconnected" | "loading";

interface ServiceContextType {
    status: ServiceStatus;
    pingService: () => Promise<void>;
}

const ServiceContext = createContext<ServiceContextType | undefined>(undefined);

export const ServiceProvider = ({ children }: { children: ReactNode }) => {
    const [status, setStatus] = useState<ServiceStatus>("loading");
    const timer = useRef<ReturnType<typeof setInterval> | null>(null);

    const pingService = async () => {
        try {
            const response = await fetch('http://localhost:8214/ping', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                const data = await response.text();
                if (data.toLowerCase().includes('pong') || response.status === 200) {
                    setStatus("connected");
                } else {
                    setStatus("disconnected");
                }
            } else {
                setStatus("disconnected");
            }
        } catch (error) {
            setStatus("disconnected");
        }
    };

    useEffect(() => {
        pingService();
        timer.current = setInterval(pingService, 5000);

        return () => {
            if (timer.current != null) clearInterval(timer.current);
        };
    }, []);

    return (
        <ServiceContext.Provider value={{ status, pingService }}>
            {children}
        </ServiceContext.Provider>
    );
};

export const useServiceStatus = () => {
    const context = useContext(ServiceContext);
    if (context === undefined) {
        throw new Error('useServiceStatus must be used within a ServiceProvider');
    }
    return context;
};
