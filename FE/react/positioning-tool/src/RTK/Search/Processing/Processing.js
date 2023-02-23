import React from "react";
import "../../../dist/css/zui.css"
import { useReq } from "../ReqProvider";

export default function Processing() {
    const {rtkProcessingState} = useReq();
    let curstate = rtkProcessingState * 100 + '%'
    return (
        <>
            <span className="label">PROCESSING</span>
            <div className="progress">
                <div className="progress-bar" role="progressbar" aria-valuenow="40" aria-valuemin="0" aria-valuemax="100" style={{width: curstate}}>
                    <span className="sr-only">curstate Complete (success)</span>
                </div>
            </div>
            <br/>
        </>
    )
}