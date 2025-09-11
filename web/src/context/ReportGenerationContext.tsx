import { createContext, useContext, useState, ReactNode } from 'react';

type ReportGenerationContextType = {
    currentGeneratedId: string | null;
    setCurrentGeneratedId: (id: string | null) => void;
    isCurrentlyGenerated: (id: string) => boolean;
    generateReport: (id: string) => void;
    resetGeneration: () => void;
};

const ReportGenerationContext = createContext<ReportGenerationContextType | undefined>(undefined);

type ReportGenerationProviderProps = {
    children: ReactNode;
};

export const ReportGenerationProvider = ({ children }: ReportGenerationProviderProps) => {
    const [currentGeneratedId, setCurrentGeneratedId] = useState<string | null>(null);

    const isCurrentlyGenerated = (id: string) => {
        return currentGeneratedId === id;
    };

    const generateReport = (id: string) => {
        setCurrentGeneratedId(id);
    };

    const resetGeneration = () => {
        setCurrentGeneratedId(null);
    };

    return (
        <ReportGenerationContext.Provider
            value={{
                currentGeneratedId,
                setCurrentGeneratedId,
                isCurrentlyGenerated,
                generateReport,
                resetGeneration,
            }}
        >
            {children}
        </ReportGenerationContext.Provider>
    );
};

export const useReportGeneration = () => {
    const context = useContext(ReportGenerationContext);
    if (!context) {
        throw new Error('useReportGeneration must be used within a ReportGenerationProvider');
    }
    return context;
};
