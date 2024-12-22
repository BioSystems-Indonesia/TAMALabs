// in src/admin/index.tsx
import {Admin, Resource} from "react-admin";
import jsonServerProvider from "ra-data-json-server";
import {DefaultLayout} from "./layout.tsx";
import {PatientCreate, PatientEdit, PatientList, PatientShow} from "./patient";
import UserIcon from '@mui/icons-material/Person';
import ScienceIcon from '@mui/icons-material/Science';
import BiotechIcon from '@mui/icons-material/Biotech';
import {SpecimenCreate, SpecimenEdit, SpecimenList, SpecimenShow} from "./specimen";
import {WorkOrderCreate, WorkOrderEdit, WorkOrderList} from "./workOrder";

const dataProvider = jsonServerProvider(import.meta.env.VITE_BACKEND_BASE_URL);

const App = () => {
    return (<Admin dataProvider={dataProvider} layout={DefaultLayout}>
        <Resource name="work-order" list={WorkOrderList}
                  edit={WorkOrderEdit}
                  create={WorkOrderCreate}
                  hasCreate={true}
                  hasEdit={true}
                  hasShow={false}
                  icon={BiotechIcon}/>
        <Resource name="patient" list={PatientList} show={PatientShow} edit={PatientEdit} create={PatientCreate}
                  hasCreate={true}
                  hasEdit={true}
                  hasShow={true}
                  icon={UserIcon}
                  recordRepresentation={record => `#${record.id} - ${record.first_name} ${record.last_name}`}
        />
        <Resource name="Specimen" list={SpecimenList} show={SpecimenShow} edit={SpecimenEdit}
                  create={SpecimenCreate}
                  hasCreate={true}
                  hasEdit={true}
                  hasShow={true}
                  icon={ScienceIcon}
                  recordRepresentation={record => `#${record.id} - ${record.type}`}
        />
    </Admin>)
}


export default App;