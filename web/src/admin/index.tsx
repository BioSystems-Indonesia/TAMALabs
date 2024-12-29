// in src/admin/index.tsx
import { Admin, CustomRoutes, Resource } from "react-admin";
import { Route } from "react-router-dom";
import jsonServerProvider from "ra-data-json-server";
import LanIcon from '@mui/icons-material/Lan';
import { DefaultLayout } from "./layout.tsx";
import { PatientCreate, PatientEdit, PatientList, PatientShow } from "./patient";
import UserIcon from '@mui/icons-material/Person';
import ScienceIcon from '@mui/icons-material/Science';
import BiotechIcon from '@mui/icons-material/Biotech';
import AssessmentIcon from '@mui/icons-material/Assessment';
import { WorkOrderAddTest, WorkOrderCreate, WorkOrderEdit, WorkOrderList } from "./workOrder";
import { ObservationRequestList, ObservationRequestShow } from "./observationRequest";
import { SpecimenList, SpecimenShow } from "./specimen";
import ListIcon from '@mui/icons-material/List';
import { WorkOrderShow } from "./workOrder/Show.tsx";
import { DeviceCreate, DeviceEdit, DeviceList, DeviceShow } from "./device/index.tsx";
import { TestTypeList, TestTypeShow } from "./testType";
import Settings from "./settings/index.tsx";

const dataProvider = jsonServerProvider(import.meta.env.VITE_BACKEND_BASE_URL);

const App = () => {
    return (<Admin dataProvider={dataProvider} layout={DefaultLayout}>
        <CustomRoutes>
            <Route path="/settings*" element={<Settings />} />
        </CustomRoutes>
        <Resource name="work-order" list={WorkOrderList}
            edit={WorkOrderEdit}
            create={WorkOrderCreate}
            show={WorkOrderShow}
            hasCreate={true}
            hasEdit={true}
            hasShow={true}
            icon={BiotechIcon}
        >
            <Route path="/:id/add-test*" element={<WorkOrderAddTest />} />
        </Resource>
        <Resource name="result" list={SpecimenList} show={SpecimenShow}
                  hasCreate={false}
                  hasEdit={false}
                  hasShow={true}
                  icon={AssessmentIcon}
                  recordRepresentation={record => `#${record.id} - ${record.type}`}
        />
        <Resource name="patient" list={PatientList} show={PatientShow} edit={PatientEdit} create={PatientCreate}
            hasCreate={true}
            hasEdit={true}
            hasShow={true}
            icon={UserIcon}
            recordRepresentation={record => `#${record.id} - ${record.first_name} ${record.last_name}`}
        />
        <Resource name="specimen" list={SpecimenList} show={SpecimenShow}
            hasCreate={false}
            hasEdit={false}
            hasShow={true}
            icon={ScienceIcon}
            recordRepresentation={record => `#${record.id} - ${record.type}`}
        />
        <Resource name="test-type" list={TestTypeList} show={TestTypeShow}
            hasCreate={false}
            hasEdit={false}
            hasShow={true}
            icon={BiotechIcon}
            recordRepresentation={record => `#${record.id} - ${record.code}`}
        />
        <Resource name="device" list={DeviceList} show={DeviceShow} edit={DeviceEdit}
            create={DeviceCreate}
            hasCreate={true}
            hasEdit={true}
            hasShow={true}
            icon={LanIcon}
            recordRepresentation={record => `#${record.id} - ${record.name}`}
        />
        <Resource name="observation-request" list={ObservationRequestList} show={ObservationRequestShow}
            hasCreate={false}
            hasEdit={false}
            hasShow={true}
            icon={ListIcon}
            recordRepresentation={record => `#${record.id} - ${record.type}`}
        />
    </Admin>)
}


export default App;
