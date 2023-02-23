import React from "react";
import "../../dist/css/zui.css"
import ReqProvider from "./ReqProvider";
import Stdout from "./Stdout/Stdout";
import FilePath from "./FilePath/FilePath";
import FileName from "./FileName/FileName";
import LatColName from "./LatColName/LatColName";
import LonColName from "./LonColName/LonColName";
import SearchButton from "./Button/SearchButton";

export default function Search() {
    return (
        <ReqProvider>
            <FilePath></FilePath>
            <br></br>
            <FileName></FileName>
            <br></br>
            <LatColName></LatColName>
            <br></br>
            <LonColName></LonColName>
            <br></br>
            <SearchButton></SearchButton>
            <br></br>
            <Stdout></Stdout>
        </ReqProvider>
    )
}