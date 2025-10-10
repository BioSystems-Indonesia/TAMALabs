import { Routes, Route, useNavigate } from "react-router-dom";
import MyAdmin from "./admin";
import { useEffect } from "react";
import axios from "axios";
import LicensePage from "./license";

const App = () => {
    const navigate = useNavigate();

    useEffect(() => {
        axios
            .get("/api/v1/license/check", {
                baseURL: import.meta.env.VITE_BACKEND_BASE_URL,
            })
            .catch(() => {
                navigate("/license");
            });
    }, [navigate]);

    return (
        <Routes>
            <Route path="/" element={<MyAdmin />} />
            <Route path="/license" element={<LicensePage />} />
        </Routes>
    );
};

export default App;
