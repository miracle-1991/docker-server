import React, {useState} from "react";
import {useReq} from "./ReqProvider";

export default function Prefix() {
    const {req, onPrefixChange} = useReq();
    const [prefixStr, setPrefixStr] = useState(req.prefix)
    return (
        <>
            <span className="label">APP</span>
            <input type="text" className="form-control" value={prefixStr} onChange={
                event => {
                    let newV = event.target.value;
                    setPrefixStr(newV);
                    onPrefixChange(newV)
                }
            }/>
        </>
    )
}