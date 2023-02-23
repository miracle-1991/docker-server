import React, {useState} from "react";
import { useReq } from "./ReqProvider";

export default function OutputPath() {
    const {req, onOutputpathChange } = useReq();
    const [dir, setDir] = useState(req.outputpath);
    return (
        <>
            <div className="row">
                <div className="col-xs-12">
                    <span className="label">SAVE PATH</span>
                    <input type="text" className="form-control" value={dir} onChange={
                        event => {
                            let newV = event.target.value;
                            setDir(newV);
                            onOutputpathChange(newV);
                        }
                    }/>
                </div>
            </div>
            <br></br>
        </>
    )
}