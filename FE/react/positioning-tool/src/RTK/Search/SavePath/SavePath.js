import React, {useState} from "react";
import { useReq } from "../ReqProvider";

export default function SavePath() {
    const {req, onOutputPathChange } = useReq();
    const [dir, setDir] = useState(req.outputpath);
    return (
        <>
            <div className="row">
                <div className="col-xs-12">
                    <span className="label">PATH TO SAVE</span>
                    <input type="text" className="form-control" value={dir} onChange={
                        event => {
                            let newV = event.target.value;
                            setDir(newV);
                            onOutputPathChange(newV);
                        }
                    }/>
                </div>
            </div>
        </>
    )
}