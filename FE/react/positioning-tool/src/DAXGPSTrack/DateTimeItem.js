import React, {useState} from "react";
import { useReq } from "./ReqProvider";

export default function DateTimeItem({titlestr=""}) {
    const {req, onTimeChange} = useReq();
    const [tstr, setTstr] = useState(req.time);
    return (
        <div className="col-xs-12">
            <span className="label">{titlestr}</span>
            <input type="datetime-local" step="1" className="form-control" value={tstr} onChange={
                event => {
                    let newT = event.target.value.replace("T", " ");
                    setTstr(newT);
                    onTimeChange(newT);
                }
            }/>
        </div>
    )
}