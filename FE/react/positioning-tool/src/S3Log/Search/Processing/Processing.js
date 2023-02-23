import React from "react";
import "../../../dist/css/zui.css"
import { useReq } from "../ReqProvider";

export default function Processing() {
    const {processingState} = useReq();
    let curstate = processingState * 100 + '%'
    return (
        <>
            <span className="label">PROCESSING</span>
            <div className="progress">
                <div className="progress-bar" role="progressbar" aria-valuenow="40" aria-valuemin="0" aria-valuemax="100" style={{width: curstate}}>
                    <span className="sr-only">{processingState} Complete (success)</span>
                </div>
            </div>
            <br/>
        </>
    )
}