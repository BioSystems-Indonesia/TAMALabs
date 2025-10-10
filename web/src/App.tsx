import dayjs from "dayjs";
import MyAdmin from "./admin";
import './App.css'
import utc from "dayjs/plugin/utc"; 
import customParseFormat from "dayjs/plugin/customParseFormat";


dayjs.extend(utc);
dayjs.extend(customParseFormat);


const App = () => <MyAdmin />;

export default App;
