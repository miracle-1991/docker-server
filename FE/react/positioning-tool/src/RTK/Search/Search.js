import React from "react";
import ReqProvider from "./ReqProvider";
import SavePath from "./SavePath/SavePath";
import Button from "./Button/Button";
import Processing from "./Processing/Processing";
import Stdout from "./Stdout/Stdout";
import Pmode from "./PMode/Pmode";
import NavSys from "./NavSys/NavSys";
import ObsNavPath from "./ObsNavPath/ObsNavPath";

export default function Search() {
    return (
        <ReqProvider>
            <Pmode></Pmode>
            <br></br>
            <NavSys></NavSys>
            <br></br>
            <ObsNavPath></ObsNavPath>
            <SavePath></SavePath>
            <br></br>
            <Button></Button>
            <br></br>
            <Processing></Processing>
            <Stdout></Stdout>
        </ReqProvider>
    )
}