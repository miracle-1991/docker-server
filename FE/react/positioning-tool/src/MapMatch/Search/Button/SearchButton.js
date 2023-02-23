import React from "react";
import {useReq} from "../ReqProvider";

export default function SearchButton() {
    const {req, snapuri, setSnapData, setSnapLoading, setError} = useReq();
    const submit = e => {
        e.preventDefault();
        console.log("submit: " + JSON.stringify(req))
        if (!snapuri) return;
        fetch(snapuri, {
            method: "post",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(req)
        }).then(res =>res.json()).then(setSnapData).then( () => {setSnapLoading(false);}).catch(setError);
    }

    return (
        <span className="input-group-btn">
            <button className="btn btn-default" type="button" onClick={submit}>SNAP BEGIN</button>
        </span>
    )
}