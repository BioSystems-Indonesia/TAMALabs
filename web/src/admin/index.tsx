import AssessmentIcon from '@mui/icons-material/Assessment';
import BiotechIcon from '@mui/icons-material/Biotech';
import ApprovalIcon from '@mui/icons-material/Approval';
import BuildIcon from '@mui/icons-material/Build';
import LanIcon from '@mui/icons-material/Lan';
import UserIcon from '@mui/icons-material/Person';
import AdminPanelSettingsIcon from '@mui/icons-material/AdminPanelSettings';
import HistoryEduIcon from '@mui/icons-material/HistoryEdu';
import TableViewIcon from '@mui/icons-material/TableView';
import InfoIcon from '@mui/icons-material/Info';
import LicenseIcon from '@mui/icons-material/VerifiedUser';
import DashboardIcon from '@mui/icons-material/Dashboard';
import FactCheckIcon from '@mui/icons-material/FactCheck';
import jsonServerProvider from "ra-data-json-server";
import { Admin, CustomRoutes, fetchUtils, HttpError, Resource } from "react-admin";
import { Route } from "react-router-dom";
import { dateFormatter } from '../helper/format.ts';
import { ConfigEdit, ConfigList } from "./config/config.tsx";
import { DeviceCreate, DeviceEdit, DeviceList, DeviceShow } from "./device/index.tsx";
import { DefaultLayout } from "./layout.tsx";
import { PatientCreate, PatientEdit, PatientList, PatientShow } from "./patient";
import PatientResultHistory from './patient/PatientResultHistory.tsx';
import { ResultList } from "./result";
import { ResultShow } from './result/show.tsx';
import Settings from "./settings/index.tsx";
import { TestTemplateCreate, TestTemplateEdit, TestTemplateList } from './testTemplate/index.tsx';
import { TestTypeCreate, TestTypeEdit, TestTypeList, TestTypeShow } from "./testType";
import { WorkOrderCreate, WorkOrderEdit, WorkOrderList } from "./workOrder";
import { WorkOrderShow } from "./workOrder/Show.tsx";
import { radiantLightTheme, radiantDarkTheme } from './theme.tsx';
import { useAuthProvider } from '../hooks/authProvider.ts';
import { UserCreate, UserEdit, UserList, UserShow } from './User/index.tsx';
import { ErrorPayload } from '../types/errors.ts';
import { ApprovalList } from './approval/index.tsx';
import CustomLoginPage from './login/index.tsx';
import LogViewer from './logView/index.tsx';
import { AboutPage } from './about/index.tsx';
import { DashboardPage } from './dashboard/index.tsx'
import LicensePage from './license/index.tsx';
import LicenseStatusPage from './licenseStatus/index.tsx';
import { QualityControlList, QualityControlDetail, QCForm, QCEntryForm } from './qualityControl';
// import { useEffect, useState } from 'react';

// Component to handle license check before authentication
// const LicenseChecker = ({ children }: { children: React.ReactNode }) => {
//     // const [licenseStatus, setLicenseStatus] = useState<'checking' | 'valid' | 'invalid'>('checking');

//     useEffect(() => {
//         // Check license immediately when app loads (before authentication)
//         fetch(`${import.meta.env.VITE_BACKEND_BASE_URL}/license/check`, {
//             method: 'GET',
//             headers: {
//                 'Accept': 'application/json',
//             },
//         })
//             .then(response => {
//                 if (!response.ok) {
//                     throw new Error('License check failed');
//                 }
//                 return response.json();
//             })
//             .then(data => {
//                 if (data.valid) {
//                     setLicenseStatus('valid');
//                 } else {
//                     setLicenseStatus('invalid');
//                 }
//             })
//             .catch(() => {
//                 // If license check fails, assume invalid
//                 setLicenseStatus('invalid');
//             });
//     }, []);

//     // Show loading while checking license
//     if (licenseStatus === 'checking') {
//         return (
//             <div style={{
//                 display: 'flex',
//                 justifyContent: 'center',
//                 alignItems: 'center',
//                 height: '100vh',
//                 fontSize: '18px'
//             }}>
//                 Checking license...
//             </div>
//         );
//     }

//     // Show license page if invalid
//     if (licenseStatus === 'invalid') {
//         return <LicensePage />;
//     }

//     // Render children if license is valid
//     return <>{children}</>;
// };


const httpClient = async (url: string, options?: fetchUtils.Options) => {
    if (!options) {
        options = {};
    }

    if (!options.headers) {
        options.headers = new Headers({ Accept: 'application/json' });
    }

    // Remove Authorization header logic since we're using cookies now
    // Cookie will be sent automatically with credentials: 'include'

    const requestHeaders = fetchUtils.createHeadersFromOptions(options);

    return fetch(url, {
        ...options,
        headers: requestHeaders,
        credentials: 'include' // This ensures cookies are sent with every request
    })
        .then(response =>
            response.text().then(text => ({
                status: response.status,
                statusText: response.statusText,
                headers: response.headers,
                body: text,
            }))
        )
        .then(({ status, statusText, headers, body }) => {
            let json;
            try {
                json = JSON.parse(body);
            } catch (e) {
                // not json, no big deal
            }
            if (status < 200 || status >= 300) {
                const errorPayload = json as ErrorPayload;
                return Promise.reject(
                    new HttpError(
                        (errorPayload?.error) || statusText,
                        status,
                        errorPayload
                    )
                );
            }
            return Promise.resolve({ status, headers, body, json });
        });

};

const dataProvider = jsonServerProvider(import.meta.env.VITE_BACKEND_BASE_URL, httpClient);

const App = () => {
    return (
        // <LicenseChecker>
        <Admin
            dataProvider={dataProvider}
            layout={DefaultLayout}
            theme={radiantLightTheme}
            darkTheme={radiantDarkTheme}
            authProvider={useAuthProvider()}
            loginPage={CustomLoginPage}
        >
            {permissions => (
                <>
                    {permissions === "Admin" && (
                        <>
                            <CustomRoutes>
                                <Route path="/settings/*" element={<Settings />} />
                            </CustomRoutes>
                            <CustomRoutes>
                                <Route path="/logs" element={<LogViewer />} />
                            </CustomRoutes>
                            <CustomRoutes>
                                <Route path="/license" element={<LicensePage />} />
                            </CustomRoutes>
                            <CustomRoutes>
                                <Route path="/license-status" element={<LicenseStatusPage />} />
                            </CustomRoutes>
                        </>
                    )}
                    <Resource name="dashboard"
                        list={DashboardPage}
                        options={{ label: "Dashboard" }}
                        icon={DashboardIcon}
                    >

                    </Resource>
                    <Resource
                        name="work-order"
                        list={WorkOrderList}
                        create={WorkOrderCreate}
                        show={WorkOrderShow}
                        hasCreate={true}
                        edit={WorkOrderEdit}
                        hasShow={true}
                        icon={BiotechIcon}
                        options={{
                            label: "Lab Request"
                        }}
                        recordRepresentation={record => `#${record.id} - ${dateFormatter(record.created_at)}`}
                    >
                        <Route path="/:id/show/device/create" element={<DeviceCreate />} />
                    </Resource>

                    <Resource name="result" list={ResultList} show={ResultShow}
                        hasCreate={false}
                        hasEdit={false}
                        hasShow={true}
                        icon={AssessmentIcon}
                        recordRepresentation={record => `#${record.barcode}}`}
                    />
                    <Resource
                        name="quality-control"
                        list={QualityControlList}
                        hasCreate={false}
                        hasEdit={false}
                        hasShow={false}
                        icon={FactCheckIcon}
                        options={{
                            label: "Quality Control"
                        }}
                    >
                        <Route path="/:id" element={<QualityControlDetail />} />
                        <Route path="/:deviceId/parameter/:testTypeId" element={<QCForm />} />
                        <Route path="/:deviceId/parameter/:testTypeId/entry/new" element={<QCEntryForm />} />
                    </Resource>
                    <Resource name="approval" list={ApprovalList} show={ResultShow}
                        hasCreate={false}
                        hasEdit={false}
                        hasShow={true}
                        icon={ApprovalIcon}
                        recordRepresentation={record => `#${record.barcode}}`}
                    />

                    <Resource name="patient" list={PatientList} show={PatientShow} edit={PatientEdit} create={PatientCreate}
                        hasCreate={true}
                        hasEdit={true}
                        hasShow={true}
                        icon={UserIcon}
                        recordRepresentation={record => `#${record.id} - ${record.first_name} ${record.last_name}`} >
                        <Route path="/:id/result/history" element={<PatientResultHistory />} />
                    </Resource>
                    <Resource name="patient-history"
                        list={PatientResultHistory}
                        hasCreate={false}
                        hasEdit={false}
                        hasShow={false}
                        icon={HistoryEduIcon}
                        options={{
                            label: "Patient History"
                        }}
                    >
                    </Resource>
                    <Resource name="test-type" list={TestTypeList} show={TestTypeShow}
                        {...(permissions !== "Analyzer" ? {
                            create: TestTypeCreate,
                            edit: TestTypeEdit,
                            hasCreate: true,
                            hasEdit: true,
                        } : {
                            hasCreate: false,
                            hasEdit: false,
                        })}
                        icon={BiotechIcon}
                        recordRepresentation={record => `#${record.id} - ${record.code}`}
                        options={{
                            label: "Test Type"
                        }}
                    />
                    <Resource name="test-template" list={TestTemplateList}
                        create={TestTemplateCreate}
                        edit={TestTemplateEdit}
                        hasCreate={true}
                        hasEdit={true}
                        hasShow={false}
                        icon={TableViewIcon}
                        recordRepresentation={record => `${record.name}`}
                        options={{
                            label: "Test Template"
                        }}
                    />
                    {permissions === "Admin" && (
                        <Resource name="device" list={DeviceList} show={DeviceShow} edit={DeviceEdit}
                            create={DeviceCreate}
                            hasCreate={true}
                            hasEdit={true}
                            hasShow={true}
                            icon={LanIcon}
                            recordRepresentation={record => `#${record.id} - ${record.name}`}
                        />
                    )}

                    <Resource name="user" list={UserList} show={UserShow}
                        {...(permissions === "Admin" ? {
                            create: UserCreate,
                            edit: UserEdit,
                            hasCreate: true,
                            hasEdit: true,
                        } : {
                            hasCreate: false,
                            hasEdit: false,
                        })}
                        icon={AdminPanelSettingsIcon}
                        recordRepresentation={record => `#${record.id} - ${record.fullname}`}
                    />

                    {permissions === "Admin" && (
                        <Resource name="config" list={ConfigList} edit={ConfigEdit}
                            hasCreate={false}
                            hasEdit={true}
                            icon={BuildIcon}
                            recordRepresentation={record => `#${record.id} - ${record.type}`}
                        />
                    )}

                    <Resource
                        name="about"
                        list={() => <AboutPage />}
                        hasCreate={false}
                        hasEdit={false}
                        hasShow={false}
                        icon={InfoIcon}
                        options={{
                            label: "About Us"
                        }}
                    />

                    {permissions === "Admin" && (
                        <Resource
                            name="license-status"
                            list={() => <LicenseStatusPage />}
                            hasCreate={false}
                            hasEdit={false}
                            hasShow={false}
                            icon={LicenseIcon}
                            options={{
                                label: "License Status"
                            }}
                        />
                    )}
                </>
            )}
        </Admin>
        // </LicenseChecker>
    )
}

export default App;
