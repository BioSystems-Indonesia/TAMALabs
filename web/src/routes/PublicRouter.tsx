import { BrowserRouter, Routes, Route } from 'react-router-dom';
import AppAdmin from '../admin';
import LicensePage from '../license';

export default function PublicRouter() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/license" element={<LicensePage />} />
                <Route path="/*" element={<AppAdmin />} />
            </Routes>
        </BrowserRouter>
    );
}
