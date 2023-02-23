import React from "react";
import "../dist/css/zui.css"
import DateTimeItem from "./DateTimeItem";

export default function DateTime() {
    return (
        <>
            <div className="row">
                <DateTimeItem titlestr="TIME OF DAX"></DateTimeItem>
            </div>
            <br></br>
        </>
    )
}