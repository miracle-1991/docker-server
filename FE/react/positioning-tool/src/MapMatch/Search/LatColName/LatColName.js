import React, {useState} from "react";
import { useReq } from "../ReqProvider";

export default function LatColName() {
    const {req, onLatitudeColumnNameChange } = useReq();
    const [name, setName] = useState(req.latitude_column_name);
    return (
        <>
            <div className="row">
                <div className="col-xs-12">
                    <span className="label">NAME OF LAT COLUMN</span>
                    <input type="text" className="form-control" value={name} onChange={
                        event => {
                            let newV = event.target.value;
                            setName(newV);
                            onLatitudeColumnNameChange(newV);
                        }
                    }/>
                </div>
            </div>
        </>
    )
}