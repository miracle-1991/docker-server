import React, {useState} from "react";
import { useReq } from "./ReqProvider";

export default function ForwardBackward() {
    const {req, onDriveridChange, onForwardChange, onBackwardChange } = useReq();
    const [driverid, setDriverid] = useState(req.driverid);
    const [forward, setForward] = useState(req.forward);
    const [backward, setBackward] = useState(req.backward);
    return (
        <>
            <div className="row">
                <div className="col-xs-4">
                    <span className="label">DAX ID</span>
                    <input type="text" className="form-control" value={driverid} onChange={
                        event => {
                            let newV = event.target.value;
                            setDriverid(Number(newV));
                            onDriveridChange(Number(newV));
                        }
                    }/>
                </div>
                <div className="col-xs-4">
                    <span className="label">Forward</span>
                    <input type="text" className="form-control" value={forward} onChange={
                        event => {
                            let newV = event.target.value;
                            setForward(Number(newV));
                            onForwardChange(Number(newV));
                        }
                    }/>
                </div>
                <div className="col-xs-4">
                    <span className="label">Forward</span>
                    <input type="text" className="form-control" value={backward} onChange={
                        event => {
                            let newV = event.target.value;
                            setBackward(Number(newV));
                            onBackwardChange(Number(newV));
                        }
                    }/>
                </div>
            </div>
            <br></br>
        </>
    )
}