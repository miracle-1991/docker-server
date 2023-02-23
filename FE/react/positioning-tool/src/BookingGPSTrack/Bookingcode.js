import React, {useState} from "react";
import { useReq } from "./ReqProvider";

export default function Bookingcode() {
    const {req, onBookingCodeChange } = useReq();
    const [bookingcode, setBookingcode] = useState(req.bookingcode);
    return (
        <>
            <div className="row">
                <div className="col-xs-4">
                    <span className="label">BOOKING CODE</span>
                    <input type="text" className="form-control" value={bookingcode} onChange={
                        event => {
                            let newV = event.target.value;
                            setBookingcode(newV);
                            onBookingCodeChange(newV);
                        }
                    }/>
                </div>
            </div>
            <br></br>
        </>
    )
}