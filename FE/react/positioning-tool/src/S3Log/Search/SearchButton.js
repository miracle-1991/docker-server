import React from "react";
import {useReq} from "./ReqProvider";

export default function SearchButton() {
    const {req, uri, setData, setLoading, setError, onDataLoading} = useReq();
    const submit = e => {
        e.preventDefault();
        console.log("submit: " + JSON.stringify(req))
        if (!uri) return;
        onDataLoading();
        fetch(uri, {
            method: "post",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(req)
        }).then(res =>res.json()).then(setData).then( () => {setLoading(false);}).catch(setError);
    }

    return (
        <span className="input-group-btn">
            <button className="btn btn-default" type="button" onClick={submit}>SEARCH</button>
        </span>
    )
}