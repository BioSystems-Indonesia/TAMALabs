import { BrowserRouter, Routes, Route } from 'react-router-dom';
import MyAdmin from "./admin";
import { DashboardWindow } from './admin/dashboard/window';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { radiantLightTheme } from './admin/theme';

const App = () => {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/dashboard-window" element={
                    <ThemeProvider theme={radiantLightTheme}>
                        <CssBaseline />
                        <DashboardWindow />
                    </ThemeProvider>
                } />
                <Route path="/*" element={<MyAdmin />} />
            </Routes>
        </BrowserRouter>
    );
};

export default App;