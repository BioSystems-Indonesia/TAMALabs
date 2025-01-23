// in src/admin/index.tsx
import AssessmentIcon from '@mui/icons-material/Assessment';
import BiotechIcon from '@mui/icons-material/Biotech';
import BuildIcon from '@mui/icons-material/Build';
import LanIcon from '@mui/icons-material/Lan';
import UserIcon from '@mui/icons-material/Person';
import TableViewIcon from '@mui/icons-material/TableView';
import jsonServerProvider from "ra-data-json-server";
import { Admin, CustomRoutes, Resource } from "react-admin";
import { Route } from "react-router-dom";
import { ConfigEdit, ConfigList } from "./config/config.tsx";
import { DeviceCreate, DeviceEdit, DeviceList, DeviceShow } from "./device/index.tsx";
import { DefaultLayout } from "./layout.tsx";
import { PatientCreate, PatientEdit, PatientList, PatientShow } from "./patient";
import { ResultList } from "./result";
import Settings from "./settings/index.tsx";
import { TestTypeCreate, TestTypeEdit, TestTypeList} from "./testType";
import { WorkOrderAddTest, WorkOrderCreate, WorkOrderList } from "./workOrder";
import { WorkOrderShow } from "./workOrder/Show.tsx";
import { TestTemplateCreate, TestTemplateEdit, TestTemplateList} from './testTemplate/index.tsx';

const dataProvider = jsonServerProvider(import.meta.env.VITE_BACKEND_BASE_URL);

const App = () => {
    return (<Admin dataProvider={dataProvider} layout={DefaultLayout}>
        <CustomRoutes>
            <Route path="/settings*" element={<Settings />} />
        </CustomRoutes>
        <Resource name="work-order" list={WorkOrderList}
            create={WorkOrderCreate}
            show={WorkOrderShow}
            hasCreate={true}
            hasShow={true}
            icon={BiotechIcon}
        >
            <Route path="/:id/show/add-test*" element={<WorkOrderAddTest />} />
            <Route path="/:id/show/patient/create" element={<PatientCreate />} />
            <Route path="/:id/show/device/create" element={<DeviceCreate />} />
        </Resource>

        <Resource name="result" list={ResultList}
            hasCreate={false}
            hasEdit={true}
            hasShow={true}
            icon={AssessmentIcon}
            recordRepresentation={record => `#${record.barcode}}`}
        />
        <Resource name="patient" list={PatientList} show={PatientShow} edit={PatientEdit} create={PatientCreate}
            hasCreate={true}
            hasEdit={true}
            hasShow={true}
            icon={UserIcon}
            recordRepresentation={record => `#${record.id} - ${record.first_name} ${record.last_name}`}
        />
        {/* <Resource name="specimen" list={SpecimenList} show={SpecimenShow}
            hasCreate={false}
            hasEdit={false}
            hasShow={true}
            icon={ScienceIcon}
            recordRepresentation={record => `#${record.id} - ${record.type}`}
        />
        <Resource name="observation-request" list={ObservationRequestList} show={ObservationRequestShow}
            hasCreate={false}
            hasEdit={false}
            hasShow={true}
            icon={ListIcon}
            recordRepresentation={record => `#${record.id} - ${record.type}`}
        /> */}
        <Resource name="test-type" list={TestTypeList} 
            create={TestTypeCreate}
            edit={TestTypeEdit}
            hasCreate={true}
            hasEdit={true}
            icon={BiotechIcon}
            recordRepresentation={record => `#${record.id} - ${record.code}`}
        />
        <Resource name="test-template" list={TestTemplateList} 
            create={TestTemplateCreate}
            edit={TestTemplateEdit}
            hasCreate={true}
            hasEdit={true}
            hasShow={false}
            icon={TableViewIcon}
            recordRepresentation={record => `${record.name}`}
        />
        <Resource name="device" list={DeviceList} show={DeviceShow} edit={DeviceEdit}
            create={DeviceCreate}
            hasCreate={true}
            hasEdit={true}
            hasShow={true}
            icon={LanIcon}
            recordRepresentation={record => `#${record.id} - ${record.name}`}
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
