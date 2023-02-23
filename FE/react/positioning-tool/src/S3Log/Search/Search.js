import React from "react";
import "../../dist/css/zui.css"
import SearchButton from "./SearchButton";
import SearchFilter from "./SearchFilter";
import DateTime from "./DateTimeText/DateTime";
import SavePath from "./SavePath/SavePath";
import ReqProvider from "./ReqProvider";
import Processing from "./Processing/Processing";
import Stdout from "./Stdout/Stdout";

export default function Search() {
    return (
        <ReqProvider>
            <DateTime></DateTime>
            <SavePath></SavePath>
            <div className="input-group">
                <SearchButton></SearchButton>
                <SearchFilter></SearchFilter>
            </div>
            <br></br>
            <Processing></Processing>
            <Stdout></Stdout>
        </ReqProvider>
    )
}