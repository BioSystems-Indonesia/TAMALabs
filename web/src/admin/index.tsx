// in src/admin/index.tsx
import AssessmentIcon from '@mui/icons-material/Assessment';
import BiotechIcon from '@mui/icons-material/Biotech';
import ApprovalIcon from '@mui/icons-material/Approval';
import BuildIcon from '@mui/icons-material/Build';
import LanIcon from '@mui/icons-material/Lan';
import UserIcon from '@mui/icons-material/Person';
import AdminPanelSettingsIcon from '@mui/icons-material/AdminPanelSettings';
import TableViewIcon from '@mui/icons-material/TableView';
import jsonServerProvider from "ra-data-json-server";
import { Admin, CustomRoutes, fetchUtils, HttpError, Resource } from "react-admin";
import { Route } from "react-router-dom";
import { dateFormatter } from '../helper/format.ts';
import { ConfigEdit, ConfigList } from "./config/config.tsx";
import { DeviceCreate, DeviceEdit, DeviceList, DeviceShow } from "./device/index.tsx";
import { DefaultLayout } from "./layout.tsx";
import { PatientCreate, PatientEdit, PatientList, PatientShow } from "./patient";
import { ResultList } from "./result";
import { ResultShow } from './result/show.tsx';
import Settings from "./settings/index.tsx";
import { TestTemplateCreate, TestTemplateEdit, TestTemplateList } from './testTemplate/index.tsx';
import { TestTypeCreate, TestTypeEdit, TestTypeList } from "./testType";
import { WorkOrderCreate, WorkOrderEdit, WorkOrderList } from "./workOrder";
import { WorkOrderShow } from "./workOrder/Show.tsx";
import { radiantLightTheme, radiantDarkTheme } from './theme.tsx';
import { useAuthProvider } from '../hooks/authProvider.ts';
import { LOCAL_STORAGE_ACCESS_TOKEN } from '../types/constant.ts';
import { UserCreate, UserEdit, UserList, UserShow } from './User/index.tsx';
import { ErrorPayload } from '../types/errors.ts';
import { ApprovalList } from './approval/index.tsx';
import CustomLoginPage from './login/index.tsx';
import LogViewer from './logView/index.tsx';

const httpClient = async (url: string, options?: fetchUtils.Options) => {
    if (!options) {
        options = {};
    }

    if (!options.headers) {
        options.headers = new Headers({ Accept: 'application/json' });
    }


    const accessToken = localStorage.getItem(LOCAL_STORAGE_ACCESS_TOKEN)
    if (accessToken) {
        //@ts-ignore
        options.headers.set('Authorization', `Bearer ${accessToken}`);
    }

    const requestHeaders = fetchUtils.createHeadersFromOptions(options);

    return fetch(url, { ...options, headers: requestHeaders })
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
    return (<Admin
        dataProvider={dataProvider}
        layout={DefaultLayout}
        theme={radiantLightTheme}
        darkTheme={radiantDarkTheme}
        authProvider={useAuthProvider()}
        loginPage={CustomLoginPage}
    >
        <CustomRoutes>
            <Route path="/settings/*" element={<Settings />} />
        </CustomRoutes>
        <CustomRoutes>
            <Route path="/logs" element={<LogViewer />} />
        </CustomRoutes>
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
            recordRepresentation={record => `#${record.id} - ${record.first_name} ${record.last_name}`}
        />
        <Resource name="test-type" list={TestTypeList}
            create={TestTypeCreate}
            edit={TestTypeEdit}
            hasCreate={true}
            hasEdit={true}
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
        <Resource name="device" list={DeviceList} show={DeviceShow} edit={DeviceEdit}
            create={DeviceCreate}
            hasCreate={true}
            hasEdit={true}
            hasShow={true}
            icon={LanIcon}
            recordRepresentation={record => `#${record.id} - ${record.name}`}
        />
        <Resource name="user" list={UserList} show={UserShow} edit={UserEdit}
            create={UserCreate}
            hasCreate={true}
            hasEdit={true}
            hasShow={true}
            icon={AdminPanelSettingsIcon}
            recordRepresentation={record => `#${record.id} - ${record.fullname}`}
        />
        <Resource name="config" list={ConfigList} edit={ConfigEdit}
            hasCreate={false}
            hasEdit={true}
            icon={BuildIcon}
            recordRepresentation={record => `#${record.id} - ${record.type}`}
        />
    </Admin>)
}


export default App;
