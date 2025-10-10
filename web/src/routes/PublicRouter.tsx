import { BrowserRouter, Routes, Route } from 'react-router-dom';
import AppAdmin from '../admin';
import LicensePage from '../license';
import CustomLoginPage from '../admin/login';

export default function PublicRouter() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/#/license" element={<LicensePage />} />
                <Route path="/#/login" element={<CustomLoginPage />} />
                <Route path="/*" element={<AppAdmin />} />
            </Routes>
        </BrowserRouter>
    );
}
