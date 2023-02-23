import React from "react";
import {useReq} from "../ReqProvider";

export default function Stdout() {
    const {snapData} = useReq();
    return (
        <>
            <div className="panel">
                <div className="panel-heading">
                    RUNNING LOG
                </div>
                <div className="panel-body">
                    <pre>{JSON.stringify(snapData, null, 2)}</pre>
                </div>
            </div>
        </>
    )
}