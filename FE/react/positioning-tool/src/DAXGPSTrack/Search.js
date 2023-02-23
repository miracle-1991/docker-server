import React from "react";
import "../dist/css/zui.css"
import SearchButton from "./SearchButton";
import ReqProvider from "./ReqProvider";
import DateTime from "./DateTime";
import ForwardBackward from "./ForwardBackward";
import OutputPath from "./OutputPath";
import Stdout from "./Stdout";

export default function Search() {
    return (
        <ReqProvider>
            <DateTime></DateTime>
            <br></br>
            <ForwardBackward></ForwardBackward>
            <br></br>
            <OutputPath></OutputPath>
            <br></br>
            <SearchButton></SearchButton>
            <br></br>
            <Stdout></Stdout>
        </ReqProvider>
    )
}