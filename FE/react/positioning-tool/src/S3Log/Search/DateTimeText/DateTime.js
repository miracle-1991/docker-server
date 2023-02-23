import React from "react";
import "../../../dist/css/zui.css"
import DateTimeItem from "./DateTimeItem";

export default function DateTime() {
    return (
        <>
            <div className="row">
                <DateTimeItem titlestr="START TIME" timetype="start"></DateTimeItem>
                <DateTimeItem titlestr="END TIME" timetype="end"></DateTimeItem>
            </div>
            <br></br>
        </>
    )
}