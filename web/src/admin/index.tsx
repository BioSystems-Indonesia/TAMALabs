// in src/admin/index.tsx
import {Admin, EditGuesser, Resource, ShowGuesser} from "react-admin";
import jsonServerProvider from "ra-data-json-server";
import {DefaultLayout} from "./layout.tsx";
import {PatientCreate, PatientEdit, PatientList} from "./patient";
import UserIcon from '@mui/icons-material/Person';

const dataProvider = jsonServerProvider(import.meta.env.VITE_BACKEND_BASE_URL);

const App = () => {

    console.log(import.meta.env.VITE_BACKEND_BASE_URL)
    return (<Admin dataProvider={dataProvider} layout={DefaultLayout}>
        <Resource name="patient" list={PatientList} show={ShowGuesser} edit={PatientEdit} create={PatientCreate}
                  hasCreate={true}
                  hasEdit={true}
                  hasShow={true}
                  icon={UserIcon}/>
    </Admin>)
}


export default App;