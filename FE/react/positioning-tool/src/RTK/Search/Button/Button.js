import React from "react";
import {useReq} from "../ReqProvider";

export default function Button() {
    const {req, rtkuri, setRtkData, setRtkLoading, setError, onRTKDataLoading} = useReq();
    const submit = e => {
        e.preventDefault();
        console.log("submit: " + JSON.stringify(req))
        if (!rtkuri) return;
        onRTKDataLoading();
        fetch(rtkuri, {
            method: "post",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(req)
        }).then(res =>res.json()).then(setRtkData).then(setRtkLoading(false)).catch(setError);
    }

    return (
        <span className="input-group-btn">
            <button className="btn btn-default" type="button" onClick={submit}>RUN RTK OFFLINE</button>
        </span>
    )
}