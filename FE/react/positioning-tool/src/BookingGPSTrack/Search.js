import React from "react";
import "../dist/css/zui.css"
import SearchButton from "./SearchButton";
import ReqProvider from "./ReqProvider";
import DateTime from "./DateTime";
import Bookingcode from "./Bookingcode";
import OutputPath from "./OutputPath";
import Stdout from "./Stdout";

export default function Search() {
    return (
        <ReqProvider>
            <Bookingcode></Bookingcode>
            <br></br>
            <DateTime></DateTime>
            <br></br>
            <OutputPath></OutputPath>
            <br></br>
            <SearchButton></SearchButton>
            <br></br>
            <Stdout></Stdout>
        </ReqProvider>
    )
}