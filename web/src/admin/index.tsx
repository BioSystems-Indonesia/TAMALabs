// in src/admin/index.tsx
import {Admin, Resource} from "react-admin";
import jsonServerProvider from "ra-data-json-server";
import {DefaultLayout} from "./layout.tsx";
import {PatientCreate, PatientEdit, PatientList, PatientShow} from "./patient";
import UserIcon from '@mui/icons-material/Person';
import ScienceIcon from '@mui/icons-material/Science';
import BiotechIcon from '@mui/icons-material/Biotech';
import {SpecimentCreate, SpecimentEdit, SpecimentList, SpecimentShow} from "./speciment";
import {WorkOrderCreate, WorkOrderEdit, WorkOrderList, WorkOrderShow} from "./workOrder";

const dataProvider = jsonServerProvider(import.meta.env.VITE_BACKEND_BASE_URL);

const App = () => {
    return (<Admin dataProvider={dataProvider} layout={DefaultLayout}>
        <Resource name="patient" list={PatientList} show={PatientShow} edit={PatientEdit} create={PatientCreate}
                  hasCreate={true}
                  hasEdit={true}
                  hasShow={true}
                  icon={UserIcon}
                  recordRepresentation={record => `#${record.id} - ${record.first_name} ${record.last_name}`}
        />
        <Resource name="speciment" list={SpecimentList} show={SpecimentShow} edit={SpecimentEdit}
                  create={SpecimentCreate}
                  hasCreate={true}
                  hasEdit={true}
                  hasShow={true}
                  icon={ScienceIcon}
                  recordRepresentation={record => `#${record.id} - ${record.description}`}
        />
        <Resource name="work-order" list={WorkOrderList} show={WorkOrderShow} edit={WorkOrderEdit}
                  create={WorkOrderCreate}
                  hasCreate={true}
                  hasEdit={true}
                  hasShow={true}
                  icon={BiotechIcon}/>
    </Admin>)
}


export default App;