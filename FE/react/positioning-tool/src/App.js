import './App.css';
import S3LOG from "./S3Log/S3LOG";
import MapMatch from "./MapMatch/MapMatch";
import {RTK} from "./RTK/RTK";
import {useRoutes} from "react-router-dom";
import GPSTrack from "./DAXGPSTrack/GPSTrack";
import BookingGPSTrack from "./BookingGPSTrack/GPSTrack";

function App() {
    let element = useRoutes([
        {path:  "/", element: <GPSTrack />},
        {path: "/booking", element: <BookingGPSTrack />},
        {path: "/rtk", element: <RTK />},
        {path: "/snap", element: <MapMatch/>},
        {path: "/s3log", element: <S3LOG/>}
    ]);
  return element;
}

export default App;
